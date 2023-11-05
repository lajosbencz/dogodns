package dns

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/lajosbencz/dogodns/pkg/config"
)

const perPage = 50

type doUpdater struct {
	doClient   *godo.Client
	tldList    []string
	recordList map[string][]godo.DomainRecord
}

func (u *doUpdater) Pull(record *Record) error {
	var domainRecord *godo.DomainRecord
	var recordIndex = -1
	for i, r := range u.recordList[record.GetTld()] {
		if r.Name == record.GetSub() {
			domainRecord = &r
			recordIndex = i
			break
		}
	}
	if domainRecord == nil {
		return fmt.Errorf("domain record does not exist in local cache: %s", record.Name)
	}
	domainRecord, _, err := u.doClient.Domains.Record(context.TODO(), record.GetTld(), domainRecord.ID)
	if err != nil {
		return err
	}
	u.recordList[record.GetTld()][recordIndex] = *domainRecord
	return u.Get(record)
}

func (u *doUpdater) Get(record *Record) error {
	has := u.Has(record)
	if !has {
		return fmt.Errorf("no such A record: %s", record.Name)
	}
	for _, r := range u.recordList[record.GetTld()] {
		if r.Name == record.GetSub() {
			record.Data = r.Data
			record.TTL = r.TTL
			record.Priority = r.Priority
			return nil
		}
	}
	return fmt.Errorf("internal error, cached records do not contain %s", record.Name)
}

func (u *doUpdater) Set(record *Record) (err error) {
	ctx := context.TODO()
	var res *godo.Response
	var existingRecord *godo.DomainRecord
	for _, r := range u.recordList[record.GetTld()] {
		if r.Name == record.GetSub() {
			existingRecord = &r
			break
		}
	}
	req := &godo.DomainRecordEditRequest{
		Type:     "A",
		Name:     record.GetSub(),
		Data:     record.Data,
		TTL:      record.TTL,
		Priority: record.Priority,
	}
	if existingRecord != nil {
		_, res, err = u.doClient.Domains.EditRecord(ctx, record.GetTld(), existingRecord.ID, req)
	} else {
		var r *godo.DomainRecord
		r, res, err = u.doClient.Domains.CreateRecord(ctx, record.GetTld(), req)
		u.recordList[record.GetTld()] = append(u.recordList[record.GetTld()], *r)
	}
	if err != nil && res.StatusCode >= 500 {
		err = ErrDOInternal
	}
	return
}

func (u *doUpdater) Has(record *Record) bool {
	for _, v := range u.recordList[record.GetTld()] {
		if v.Name == record.GetSub() {
			return true
		}
	}
	return false
}

func (u *doUpdater) pullTldList() error {
	u.tldList = []string{}
	page := 1
	for {
		var resp *godo.Response
		domains, resp, err := u.doClient.Domains.List(context.TODO(), &godo.ListOptions{
			Page:    page,
			PerPage: perPage,
		})
		if err != nil {
			if resp != nil && resp.StatusCode >= 500 {
				return ErrDOInternal
			}
			return err
		}
		for _, d := range domains {
			u.tldList = append(u.tldList, d.Name)
		}
		if resp.Links.Pages == nil || resp.Links.Pages.Next == "" {
			break
		}
		page++
	}
	return nil
}

func (u *doUpdater) pullRecordList() error {
	u.recordList = map[string][]godo.DomainRecord{}
	for _, tld := range u.tldList {
		page := 1
		for {
			records, resp, err := u.doClient.Domains.Records(context.TODO(), tld, &godo.ListOptions{
				Page:    page,
				PerPage: perPage,
			})
			if err != nil {
				return err
			}
			u.recordList[tld] = append(u.recordList[tld], records...)
			if resp.Links.Pages == nil || resp.Links.Pages.Next == "" {
				break
			}
			page++
		}
	}
	return nil
}

func (u *doUpdater) pullInfo() error {
	if err := u.pullTldList(); err != nil {
		return err
	}
	if err := u.pullRecordList(); err != nil {
		return err
	}
	return nil
}

func DefaultUpdater(cfg config.Config) (Updater, error) {
	u := &doUpdater{}
	u.doClient = godo.NewFromToken(cfg.Token)
	if err := u.pullInfo(); err != nil {
		return nil, err
	}
	return u, nil
}
