package domain

import (
	"context"
	"time"
)

type WeibyStoreInfo struct {
	PartnerId      string    `json:"partner_id"`
	Name           string    `json:"name"`
	Note           string    `json:"note"`
	PartnerName    string    `json:"partner_name"`
	PartnerType    int       `json:"partner_type"`
	Status         int       `json:"status"`
	ActivateDate   time.Time `json:"activate_date"`
	DeactivateDate time.Time `json:"deactivate_date"`
	CreateAt       time.Time `json:"create_at"`
	UpdateAt       time.Time `json:"update_at"`
}

type WeibyRepository interface {
	GetStoreList(ctx context.Context) (result []*WeibyStoreInfo, err error)
	GetStore(ctx context.Context, pid string) (result *WeibyStoreInfo, err error)
	GetOrderList(ctx context.Context, pid, startTime, endTime string, ptype int) (result *OrderList, err error)
}

type WeibyUsecase interface {
	GetStoreList(ctx context.Context) ([]*WeibyStoreInfo, error)
	GetStore(ctx context.Context, pid string) (*WeibyStoreInfo, error)
	GetOrderList(ctx context.Context, pid, startTime, endTime string) (*OrderList, error)
}

type UberEatsOrder struct {
	Id                        string               `json:"id"`
	DisplayId                 string               `json:"display_id"`
	ExternalReferenceId       string               `json:"external_reference_id"`
	CurrentState              string               `json:"current_state"`
	Type                      string               `json:"type"`
	Brand                     string               `json:"brand"`
	Store                     UberEatsStore        `json:"store"`
	Eater                     UberEatsEater        `json:"eater"`
	Eaters                    []UberEatsEater      `json:"eaters"`
	Cart                      UberEatsCart         `json:"cart"`
	Payment                   UberEatsPayment      `json:"payment"`
	Packaging                 UberEatsPackaging    `json:"packaging"`
	PlacedAt                  string               `json:"placed_at"`
	EstimatedReadyForPickupAt string               `json:"estimated_ready_for_pickup_at"`
	Deliveries                UberEatsEatsDelivery `json:"deliveries"`
	OrderManagerClientId      string               `json:"order_manager_client_id"`
}

type UberEatsStore struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	IntegratorStoreId string `json:"integrator_store_id"`
	IntegratorBrandId string `json:"integrator_brand_id"`
	MerchantStoreId   string `json:"merchant_store_id"`
}

type UberEatsEater struct {
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Phone     string           `json:"phone"`
	PhoneCode string           `json:"phone_code"`
	Id        string           `json:"id"`
	Delivery  UberEatsDelivery `json:"delivery"`
}

type UberEatsDelivery struct {
	Location UberEatsLocation `json:"location"`
	Type     string           `json:"type"`
	Notes    string           `json:"notes"`
}

type UberEatsEatsDelivery struct {
	Id                  string          `json:"id"`
	FirstName           string          `json:"first_name"`
	Vehicle             UberEatsVehicle `json:"vehicle"`
	PictureUrl          string          `json:"picture_url"`
	EstimatedPickupTime string          `json:"estimated_pickup_time"`
	CurrentState        string          `json:"current_state"`
	Phone               string          `json:"phone"`
	PhoneCode           string          `json:"phone_code"`
}

type UberEatsVehicle struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	LicensePlate string `json:"license_plate"`
}

type UberEatsLocation struct {
	Type          string `json:"type"`
	StreetAddress string `json:"street_address"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	GooglePlaceId string `json:"google_place_id"`
	UnitNumber    string `json:"unit_number"`
	BusinessName  string `json:"business_name"`
	Title         string `json:"title"`
}

type UberEatsCart struct {
	Items               []UberEatsItem           `json:"items"`
	SpecialInstructions string                   `json:"special_instructions"`
	FulfillmentIssues   UberEatsFulfillmentIssue `json:"fulfillment_issues"`
}

type UberEatsItem struct {
	Id                     string                    `json:"id"`
	InstanceId             string                    `json:"instance_id"`
	Title                  string                    `json:"title"`
	ExternalData           string                    `json:"external_data"`
	Quantity               int                       `json:"quantity"`
	Price                  UberEatsItemPrice         `json:"price"`
	SelectedModifierGroups UberEatsModifierGroup     `json:"selected_modifier_groups"`
	SpecialRequests        UberEatsSpecialRequests   `json:"special_requests"`
	DefaultQuantity        int                       `json:"default_quantity"`
	SpecialInstructions    string                    `json:"special_instructions"`
	FulfillmentAction      UberEatsFulfillmentAction `json:"fulfillment_action"`
	EaterId                string                    `json:"eater_id"`
	TaxInfo                UberEatsTaxInfo           `json:"tax_info"`
}

type UberEatsItemPrice struct {
	UnitPrice      UberEatsMoney `json:"unit_price"`
	TotalPrice     UberEatsMoney `json:"total_price"`
	BaseUnitPrice  UberEatsMoney `json:"base_unit_price"`
	BaseTotalPrice UberEatsMoney `json:"base_total_price"`
}

type UberEatsMoney struct {
	Amount          int    `json:"amount"`
	CurrencyCode    string `json:"currency_code"`
	FormattedAmount string `json:"formatted_amount"`
}

type UberEatsModifierGroup struct {
	Id            string         `json:"id"`
	Title         string         `json:"title"`
	ExternalData  string         `json:"external_data"`
	SelectedItems []UberEatsItem `json:"selected_items"`
	RemovedItems  []UberEatsItem `json:"removed_items"`
}

type UberEatsSpecialRequests struct {
	Allergy UberEatsAllergy `json:"allergy"`
}

type UberEatsAllergy struct {
	AllergensToExclude  []UberEatsAllergen `json:"allergens_to_exclude"`
	AllergyInstructions string             `json:"allergy_instructions"`
}

type UberEatsAllergen struct {
	Type         string `json:"type"`
	FreeformText string `json:"freeform_text"`
}

type UberEatsFulfillmentAction struct {
	FulfillmentActionType string         `json:"fulfillment_action_type"`
	ItemSubstitutes       []UberEatsItem `json:"item_substitutes"`
}

type UberEatsTaxInfo struct {
	Labels []string `json:"labels"`
}

type UberEatsFulfillmentIssue struct {
	FulfillmentIssueType  string                       `json:"fulfillment_issue_type"`
	FulfillmentActionType string                       `json:"fulfillment_action_type"`
	RootItem              UberEatsItem                 `json:"root_item"`
	ItemAvailabilityInfo  UberEatsItemAvailabilityInfo `json:"item_availability_info"`
	ItemSubstitute        UberEatsItem                 `json:"item_substitute"`
}

type UberEatsItemAvailabilityInfo struct {
	ItemsRequested int `json:"items_requested"`
	ItemsAvailable int `json:"items_available"`
}

type UberEatsPayment struct {
	Charges    UberEatsCharges    `json:"charges"`
	Accounting UberEatsAccounting `json:"accounting"`
	Promotions UberEatsPromotions `json:"promotions"`
}

type UberEatsCharges struct {
	Total                UberEatsMoney `json:"total"`
	SubTotal             UberEatsMoney `json:"sub_total"`
	Tax                  UberEatsMoney `json:"tax"`
	TotalFee             UberEatsMoney `json:"total_fee"`
	TotalFeeTax          UberEatsMoney `json:"total_fee_tax"`
	BagFee               UberEatsMoney `json:"bag_fee"`
	TotalPromoApplied    UberEatsMoney `json:"total_promo_applied"`
	SubTotalPromoApplied UberEatsMoney `json:"sub_total_promo_applied"`
	TaxPromoApplied      UberEatsMoney `json:"tax_promo_applied"`
	PickAndPackFee       UberEatsMoney `json:"pick_and_pack_fee"`
	DeliveryFee          UberEatsMoney `json:"delivery_fee"`
	DeliveryFeeTax       UberEatsMoney `json:"delivery_fee_tax"`
	SmallOrderFee        UberEatsMoney `json:"small_order_fee"`
	SmallOrderFeeTax     UberEatsMoney `json:"small_order_fee_tax"`
	Tip                  UberEatsMoney `json:"tip"`
	CashAmountDue        UberEatsMoney `json:"cash_amount_due"`
}

type UberEatsAccounting struct {
	TaxRemittance UberEatsTaxRemittance
	TaxReporting  UberEatsTaxReporting
}

type UberEatsTaxRemittance struct {
	Tax              UberEatsRemittanceInfo `json:"tax"`
	TotalFeeTax      UberEatsRemittanceInfo `json:"total_fee_tax"`
	DeliveryFeeTax   UberEatsRemittanceInfo `json:"delivery_fee_tax"`
	SmallOrderFeeTax UberEatsRemittanceInfo `json:"small_order_fee_tax"`
}

type UberEatsRemittanceInfo struct {
	Uber       []UberEatsPayeeDetail `json:"uber"`
	Restaurant []UberEatsPayeeDetail `json:"restaurant"`
	Courier    []UberEatsPayeeDetail `json:"courier"`
	Eater      []UberEatsPayeeDetail `json:"eater"`
}

type UberEatsPayeeDetail struct {
	Value UberEatsMoney `json:"value"`
}

type UberEatsTaxReporting struct {
	Breakdown   UberEatsTaxBreakDown `json:"breakdown"`
	Origin      UberEatsTaxLocation  `json:"origin"`
	Destination UberEatsTaxLocation  `json:"destination"`
}

type UberEatsTaxBreakDown struct {
	Items      []UberEatsTaxBreakdownTaxInfo `json:"items"`
	Fees       []UberEatsTaxBreakdownTaxInfo `json:"fees"`
	Promotions []UberEatsTaxBreakdownTaxInfo `json:"promotions"`
}

type UberEatsTaxBreakdownTaxInfo struct {
	InstanceId  string        `json:"instance_id"`
	Type        string        `json:"type"`
	GrossAmount UberEatsMoney `json:"gross_amount"`
	NetAmount   UberEatsMoney `json:"net_amount"`
	TotalTax    UberEatsMoney `json:"total_tax"`
	Taxes       UberEatsTax   `json:"taxes"`
}

type UberEatsTax struct {
	Rate          string               `json:"rate"`
	TaxAmount     UberEatsMoney        `json:"tax_amount"`
	IsInclusive   bool                 `json:"is_inclusive"`
	Jurisdiction  UberEatsJurisdiction `json:"jurisdiction"`
	Imposition    UberEatsImposition   `json:"imposition"`
	TaxRemittance string               `json:"tax_remittance"`
}

type UberEatsJurisdiction struct {
	Level string `json:"level"`
	Name  string `json:"name"`
}

type UberEatsImposition struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

type UberEatsTaxLocation struct {
	Id          string `json:"id"`
	CountryIso2 string `json:"country_iso2"`
	PostalCode  string `json:"postal_code"`
}

type UberEatsPromotions struct {
	ExternalPromotionId     string                 `json:"external_promotion_id"`
	PromoType               string                 `json:"promo_type"`
	PromoDiscountValue      int                    `json:"promo_discount_value"`
	PromoDiscountPercentage int                    `json:"promo_discount_percentage"`
	PromoDeliveryFeeValue   int                    `json:"promo_delivery_fee_value"`
	DiscountItems           []UberEatsDiscountItem `json:"discount_items"`
}

type UberEatsDiscountItem struct {
	ExternalId            string `json:"external_id"`
	DiscountedQuantity    int    `json:"discounted_quantity"`
	DiscountAmountApplied int    `json:"discount_amount_applied"`
}

type UberEatsPackaging struct {
	DisposableItems UberEatsDisposableItems `json:"disposable_items"`
}

type UberEatsDisposableItems struct {
	ShouldInclude bool `json:"should_include"`
}

type FoodpandaOrder struct {
	Token              string                      `json:"token"`
	Code               string                      `json:"code"`
	ShortCode          string                      `json:"shortCode"`
	PreOrder           bool                        `json:"preOrder"`
	ExpiryDate         time.Time                   `json:"expiryDate"`
	CreatedAt          time.Time                   `json:"createdAt"`
	LocalInfo          FoodpandaLocalInfo          `json:"localInfo"`
	PlatformRestaurant FoodpandaPlatformRestaurant `json:"platformRestaurant"`
	Customer           FoodpandaCustomer           `json:"customer"`
	Payment            FoodpandaPayment            `json:"payment"`
	ExpeditionType     string                      `json:"expeditionType"`
	Products           []FoodpandaProduct          `json:"products"`
	Comments           FoodpandaComment            `json:"comments"`
	Discounts          FoodpandaDiscount           `json:"discounts"`
	Price              FoodpandaPrice              `json:"price"`
	WebOrder           bool                        `json:"webOrder"`
	MobileOrder        bool                        `json:"mobileOrder"`
	CorporateOrder     bool                        `json:"corporateOrder"`
	Delivery           FoodpandaDelivery           `json:"delivery"`
}

type FoodpandaLocalInfo struct {
	Platform               string `json:"platform"`
	PlatformKey            string `json:"platformKey"`
	CountryCode            string `json:"countryCode"`
	CurrencySymbol         string `json:"currencySymbol"`
	CurrencySymbolPosition string `json:"currencySymbolPosition"`
	CurrencySymbolSpaces   string `json:"currencySymbolSpaces"`
	DecimalSeparator       string `json:"decimalSeparator"`
	DecimalDigits          string `json:"decimalDigits"`
	ThousandsSeparator     string `json:"thousandsSeparator"`
	Website                string `json:"website"`
	Email                  string `json:"email"`
	Phone                  string `json:"phone"`
}

type FoodpandaPlatformRestaurant struct {
	Id string `json:"id"`
}

type FoodpandaCustomer struct {
	Id                     string `json:"id"`
	Code                   string `json:"code"`
	MobilePhone            string `json:"mobilePhone"`
	FirstName              string `json:"firstName"`
	LastName               string `json:"lastName"`
	Email                  string `json:"email"`
	MobilePhoneCountryCode string `json:"mobilePhoneCountryCode"`
}

type FoodpandaPayment struct {
	Type               string `json:"type"`
	RemoteCode         string `json:"remoteCode"`
	Status             string `json:"status"`
	RequireMoneyChange string `json:"requireMoneyChange"`
	VatName            string `json:"vatName"`
	VatId              string `json:"vatId"`
}

type FoodpandaProduct struct {
	Id               string                            `json:"id"`
	RemoteCode       string                            `json:"remoteCode"`
	Name             string                            `json:"name"`
	Description      string                            `json:"description"`
	Comment          string                            `json:"comment"`
	CategoryName     string                            `json:"categoryName"`
	Variation        FoodpandaProductVariation         `json:"variation"`
	UnitPrice        string                            `json:"unitPrice"`
	PaidPrice        string                            `json:"paidPrice"`
	DiscountAmount   string                            `json:"discountAmount"`
	Quantity         string                            `json:"quantity"`
	HalfHalf         bool                              `json:"halfHalf"`
	VatPercentage    string                            `json:"vatPercentage"`
	SelectedToppings []FoodpandaProductSelectedTopping `json:"selectedToppings"`
}

type FoodpandaProductVariation struct {
	Name string `json:"name"`
}

type FoodpandaProductSelectedTopping struct {
	RemoteCode string                            `json:"remoteCode"`
	Name       string                            `json:"name"`
	Quantity   int                               `json:"quantity"`
	Children   []FoodpandaProductSelectedTopping `json:"children"`
	Price      string                            `json:"price"`
}

type FoodpandaComment struct {
	CustomerComment string `json:"customerComment"`
	VendorComment   string `json:"vendorComment"`
}

type FoodpandaDiscount struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type FoodpandaPrice struct {
	MinimumDeliveryValue             string                 `json:"minimumDeliveryValue"`
	Comission                        string                 `json:"comission"`
	DeliveryFee                      string                 `json:"deliveryFee"`
	DeliveryFees                     []FoodpandaDeliveryFee `json:"deliveryFees"`
	ContainerCharge                  string                 `json:"containerCharge"`
	DeliveryFeeDiscount              string                 `json:"deliveryFeeDiscount"`
	ServiceFeePercent                string                 `json:"serviceFeePercent"`
	ServiceFeeTotal                  string                 `json:"serviceFeeTotal"`
	ServiceTax                       int                    `json:"serviceTax"`
	ServiceTaxValue                  int                    `json:"serviceTaxValue"`
	SubTotal                         string                 `json:"subTotal"`
	VatVisible                       bool                   `json:"vatVisible"`
	VatPercent                       string                 `json:"vatPercent"`
	VatTotal                         string                 `json:"vatTotal"`
	GrandTotal                       string                 `json:"grandTotal"`
	DiscountAmountTotal              string                 `json:"discountAmountTotal"`
	DifferenceToMinimumDeliveryValue string                 `json:"differenceToMinimumDeliveryValue"`
	PayRestaurant                    string                 `json:"payRestaurant"`
	CollectFromCustomer              string                 `json:"collectFromCustomer"`
	RiderTip                         string                 `json:"riderTip"`
}

type FoodpandaDeliveryFee struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type FoodpandaDelivery struct {
	ExpressDelivery      bool             `json:"expressDelivery"`
	ExpectedDeliveryTime time.Time        `json:"expectedDeliveryTime"`
	RiderPickupTime      string           `json:"riderPickupTime"`
	Address              FoodpandaAddress `json:"address"`
}

type FoodpandaAddress struct {
	Line1                    string  `json:"line1"`
	Line2                    string  `json:"line2"`
	Line3                    string  `json:"line3"`
	Line4                    string  `json:"line4"`
	Line5                    string  `json:"line5"`
	Street                   string  `json:"street"`
	Number                   string  `json:"number"`
	Room                     string  `json:"room"`
	FlatNumber               string  `json:"flatNumber"`
	Building                 string  `json:"building"`
	Intercom                 string  `json:"intercom"`
	Entrance                 string  `json:"entrance"`
	Structure                string  `json:"structure"`
	Floor                    string  `json:"floor"`
	District                 string  `json:"district"`
	Other                    string  `json:"other"`
	City                     string  `json:"city"`
	Postcode                 string  `json:"postcode"`
	Company                  string  `json:"company"`
	DeliveryMainArea         string  `json:"deliveryMainArea"`
	DeliveryMainAreaPostcode string  `json:"deliveryMainAreaPostcode"`
	DeliveryArea             string  `json:"deliveryArea"`
	DeliveryAreaPostcode     string  `json:"deliveryAreaPostcode"`
	DeliveryInstructions     string  `json:"deliveryInstructions"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
}

type OrderList struct {
	UberEats  []UberEatsOrder  `json:"uberEats"`
	Foodpanda []FoodpandaOrder `json:"foodpanda"`
}
