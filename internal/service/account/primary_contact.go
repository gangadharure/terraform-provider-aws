package account

import (
	"context"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/account"
	"github.com/aws/aws-sdk-go-v2/service/account/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_account_primary_contact")
func ResourcePrimaryContact() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourcePrimaryContactPut,
		ReadWithoutTimeout:   resourcePrimaryContactRead,
		UpdateWithoutTimeout: resourcePrimaryContactPut,
		DeleteWithoutTimeout: schema.NoopContext,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
			"address_line_1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_line_2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_line_3": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"city": {
				Type:     schema.TypeString,
				Required: true,
			},
			"company_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"district_or_county": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"phone_number": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\s0-9()+-]+$`), "must be a valid phone number"),
			},
			"postal_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state_or_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"website_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePrimaryContactPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AccountClient()

	id := "default"
	input := &account.PutContactInformationInput{
		ContactInformation: &types.ContactInformation{
			AddressLine1: aws.String(d.Get("address_line_1").(string)),
			City:         aws.String(d.Get("city").(string)),
			CountryCode:  aws.String(d.Get("country_code").(string)),
			FullName:     aws.String(d.Get("full_name").(string)),
			PhoneNumber:  aws.String(d.Get("phone_number").(string)),
			PostalCode:   aws.String(d.Get("postal_code").(string)),
		},
	}

	if v, ok := d.GetOk("account_id"); ok {
		id = v.(string)
		input.AccountId = aws.String(id)
	}

	if v, ok := d.GetOk("address_line_2"); ok {
		input.ContactInformation.AddressLine2 = aws.String(v.(string))
	}

	if v, ok := d.GetOk("address_line_3"); ok {
		input.ContactInformation.AddressLine3 = aws.String(v.(string))
	}

	if v, ok := d.GetOk("company_name"); ok {
		input.ContactInformation.CompanyName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("district_or_county"); ok {
		input.ContactInformation.DistrictOrCounty = aws.String(v.(string))
	}

	if v, ok := d.GetOk("state_or_region"); ok {
		input.ContactInformation.StateOrRegion = aws.String(v.(string))
	}

	if v, ok := d.GetOk("website_url"); ok {
		input.ContactInformation.WebsiteUrl = aws.String(v.(string))
	}

	_, err := conn.PutContactInformation(ctx, input)

	if err != nil {
		return diag.Errorf("creating Account Primary Contact (%s): %s", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id)
	}

	return resourcePrimaryContactRead(ctx, d, meta)
}

func resourcePrimaryContactRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).AccountClient()

	contactInformation, err := FindContactInformation(ctx, conn, d.Get("account_id").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Account Primary Contact (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading Account Primary Contact (%s): %s", d.Id(), err)
	}

	d.Set("account_id", d.Get("account_id"))
	d.Set("address_line_1", contactInformation.AddressLine1)
	d.Set("address_line_2", contactInformation.AddressLine2)
	d.Set("address_line_3", contactInformation.AddressLine3)
	d.Set("city", contactInformation.City)
	d.Set("company_name", contactInformation.CompanyName)
	d.Set("country_code", contactInformation.CountryCode)
	d.Set("district_or_county", contactInformation.DistrictOrCounty)
	d.Set("full_name", contactInformation.FullName)
	d.Set("phone_number", contactInformation.PhoneNumber)
	d.Set("postal_code", contactInformation.PostalCode)
	d.Set("state_or_region", contactInformation.StateOrRegion)
	d.Set("website_url", contactInformation.WebsiteUrl)

	return nil
}

func FindContactInformation(ctx context.Context, conn *account.Client, accountID string) (*types.ContactInformation, error) {
	input := &account.GetContactInformationInput{}
	if accountID != "" {
		input.AccountId = aws.String(accountID)
	}

	output, err := conn.GetContactInformation(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.ContactInformation == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.ContactInformation, nil
}
