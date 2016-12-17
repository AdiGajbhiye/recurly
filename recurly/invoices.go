package recurly

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type (
	// InvoicesService handles communication with theinvoices related methods
	// of the recurly API.
	InvoicesService struct {
		client *Client
	}

	// Invoice is an individual invoice for an account.
	Invoice struct {
		XMLName               xml.Name   `xml:"invoice,omitempty"`
		AccountCode           hrefString `xml:"account,omitempty"` // Read only
		Address               Address    `xml:"address,omitempty"`
		SubscriptionUUID      hrefString `xml:"subscription,omitempty"`
		OriginalInvoiceNumber hrefInt    `xml:"original_invoice,omitempty"` // Read only
		UUID                  string     `xml:"uuid,omitempty"`
		State                 string     `xml:"state,omitempty"`
		InvoiceNumberPrefix   string     `xml:"invoice_number_prefix,omitempty"`
		InvoiceNumber         int        `xml:"invoice_number,omitempty"`
		PONumber              string     `xml:"po_number,omitempty"`
		VATNumber             string     `xml:"vat_number,omitempty"`
		SubtotalInCents       int        `xml:"subtotal_in_cents,omitempty"`
		TaxInCents            int        `xml:"tax_in_cents,omitempty"`
		TotalInCents          int        `xml:"total_in_cents,omitempty"`
		Currency              string     `xml:"currency,omitempty"`
		CreatedAt             NullTime   `xml:"created_at,omitempty"`
		ClosedAt              NullTime   `xml:"closed_at,omitempty"`
		TaxType               string     `xml:"tax_type,omitempty"`
		TaxRegion             string     `xml:"tax_region,omitempty"`
		TaxRate               float64    `xml:"tax_rate,omitempty"`
		NetTerms              NullInt    `xml:"net_terms,omitempty"`
		CollectionMethod      string     `xml:"collection_method,omitempty"`
		// Redemption ? UUID is diffferent from others @todo
		LineItems    []Adjustment  `xml:"line_items>adjustment,omitempty"`
		Transactions []Transaction `xml:"transactions>transaction,omitempty"`
	}
)

const (
	// InvoiceStateOpen is an invoice state for invoices that are open, pending
	// collection.
	InvoiceStateOpen = "open"

	// InvoiceStateCollected is an invoice state for invoices that have been
	// successfully collected.
	InvoiceStateCollected = "collected"

	// InvoiceStateFailed is an invoice state for invoices that failed to collect.
	InvoiceStateFailed = "failed"

	// InvoiceStatePastDue is an invoice state for invoices where initial collection
	// failed, but Recurly is still attempting collection.
	InvoiceStatePastDue = "past_due"
)

// List returns a list of all invoices.
// https://dev.recurly.com/docs/list-invoices
func (s *InvoicesService) List(params Params) (*Response, []Invoice, error) {
	req, err := s.client.newRequest("GET", "invoices", params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName  xml.Name  `xml:"invoices"`
		Invoices []Invoice `xml:"invoice"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.Invoices, err
}

// ListAccount returns a list of all invoices for an account.
// https://dev.recurly.com/docs/list-an-accounts-invoices
func (s *InvoicesService) ListAccount(accountCode string, params Params) (*Response, []Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices", accountCode)
	req, err := s.client.newRequest("GET", action, params, nil)
	if err != nil {
		return nil, nil, err
	}

	var p struct {
		XMLName  xml.Name  `xml:"invoices"`
		Invoices []Invoice `xml:"invoice"`
	}
	resp, err := s.client.do(req, &p)

	return resp, p.Invoices, err
}

// Get returns detailed information about an invoice including line items and
// payments.
// https://dev.recurly.com/docs/lookup-invoice-details
func (s *InvoicesService) Get(invoiceNumber int) (*Response, *Invoice, error) {
	action := fmt.Sprintf("invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// GetPDF retrieves the invoice as a PDF.
// The language parameters allows you to specify a language to translate the
// invoice into. If empty, English will be used. Options: Danish, German,
// Spanish, French, Hindi, Japanese, Dutch, Portuguese, Russian, Turkish, Chinese.
// https://dev.recurly.com/docs/retrieve-a-pdf-invoice
func (s *InvoicesService) GetPDF(invoiceNumber int, language string) (*Response, *bytes.Buffer, error) {
	action := fmt.Sprintf("invoices/%d", invoiceNumber)
	req, err := s.client.newRequest("GET", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	if language == "" {
		language = "English"
	}

	req.Header.Set("Accept", "application/pdf")
	req.Header.Set("Accept-Language", language)

	var pdf bytes.Buffer
	resp, err := s.client.do(req, &pdf)

	return resp, &pdf, err
}

// Preview allows you to display the invoice details, including estimated tax,
// before you post it.
// https://dev.recurly.com/docs/post-an-invoice-invoice-pending-charges-on-an-acco
func (s *InvoicesService) Preview(accountCode string) (*Response, *Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices/preview", accountCode)
	req, err := s.client.newRequest("POST", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// Create posts an accounts pending charges to a new invoice on that account.
// When you post one-time charges to an account, these will remain pending
// until they are invoiced. An account is automatically invoiced when the
// subscription renews. However, there are times when it is appropriate to
// invoice an account before the renewal. If the subscriber has a yearly
// subscription, you might want to collect the one-time charges well before the renewal.
// https://dev.recurly.com/docs/post-an-invoice-invoice-pending-charges-on-an-acco
func (s *InvoicesService) Create(accountCode string, invoice Invoice) (*Response, *Invoice, error) {
	action := fmt.Sprintf("accounts/%s/invoices", accountCode)
	req, err := s.client.newRequest("POST", action, nil, invoice)
	if err != nil {
		return nil, nil, err
	}

	var dst Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// MarkAsPaid marks an invoice as paid successfully.
// https://dev.recurly.com/docs/mark-an-invoice-as-paid-successfully
func (s *InvoicesService) MarkAsPaid(invoiceNumber int) (*Response, *Invoice, error) {
	action := fmt.Sprintf("invoices/%d/mark_successful", invoiceNumber)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}

// MarkAsFailed marks an invoice as failed.
// https://dev.recurly.com/docs/mark-an-invoice-as-failed-collection
func (s *InvoicesService) MarkAsFailed(invoiceNumber int) (*Response, *Invoice, error) {
	action := fmt.Sprintf("invoices/%d/mark_failed", invoiceNumber)
	req, err := s.client.newRequest("PUT", action, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var dst Invoice
	resp, err := s.client.do(req, &dst)

	return resp, &dst, err
}