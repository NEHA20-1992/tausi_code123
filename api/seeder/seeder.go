package seeder

import (
	"os"
	"time"

	"github.com/NEHA20-1992/tausi_code/api/model"
	"github.com/NEHA20-1992/tausi_code/api/service"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	Raw                       = 1
	Processed                 = 2
	DtFormatSimpleDate string = "2006-01-02"
	shouldDropTables          = true
)

func dropTables(db *gorm.DB) (err error) {
	err = db.Migrator().DropTable(&model.Country{}, &model.City{}, &model.UserRole{}, &model.CompanyType{},
		&model.Company{}, &model.CompanyFile{}, &model.CompanyDataFile{}, &model.CompanyDataFileType{}, &model.User{}, &model.CustomerIdProofType{},
		&model.Model{}, &model.ModelVariable{}, &model.EnumerationItem{},
		&model.CustomerInformation{}, &model.CustomerInformationItem{})
	return
}

func Initialize(db *gorm.DB) {
	var err error
	if shouldDropTables {
		err = dropTables(db)
		if err != nil {
			panic(err)
		}
	}

	err = db.AutoMigrate(&model.Country{}, &model.City{}, &model.UserRole{}, &model.CompanyType{},
		&model.Company{}, &model.CompanyFile{}, &model.CompanyDataFile{}, &model.CompanyDataFileType{},
		&model.User{}, &model.CustomerIdProofType{}, &model.Model{}, &model.ModelVariable{}, &model.EnumerationItem{},
		&model.CustomerInformation{}, &model.CustomerInformationItem{})
	if err != nil {
		panic(err)
	}

	logoContents, err := os.ReadFile("./conf/logo.png")
	if err != nil {
		panic(err)
	}

	var userRoleService service.UserRoleService = service.GetUserRoleService(db)
	urList, err := userRoleService.GetAll()
	if err != nil {
		panic(err)
	} else if len(urList) > 0 {
		return
	}

	defaultCountry := model.Country{Name: "Tanzania, United Republic of"}
	err = loadCountries(&defaultCountry, db)
	if err != nil {
		panic(err)
	}
	err = loadCities(db)
	if err != nil {
		panic(err)
	}

	administratorCompanyType := model.CompanyType{Name: "Administrator"}
	bankCompanyType := model.CompanyType{Name: "Bank"}
	var ctList = []*model.CompanyType{
		&administratorCompanyType,
		{Name: "Mobile Network Operator"},
		&bankCompanyType,
		{Name: "Micro-Finance Institutions and Co-Operative Unions"},
		{Name: "Non-Governmental Organizations"},
		{Name: "Governmental Organizations"},
		{Name: "Financial Technology Companies"},
		{Name: "Non-Financial Companies (Private)"},
	}
	db.CreateInBatches(ctList, 100)

	var ciptList = []*model.CustomerIdProofType{
		{Name: "National Identification Number"},
		{Name: "Passport"},
		{Name: "Driving License"},
	}
	db.CreateInBatches(ciptList, 100)

	var cdftList = []*model.CompanyDataFileType{
		{Name: "Raw"},
		{Name: "Processed"},
	}
	db.CreateInBatches(cdftList, 100)

	timestamp := time.Now()

	companyOne := model.Company{
		Name:          "AfroPavo Analytics",
		CompanyTypeID: administratorCompanyType.ID,
		Address:       "2nd Floor Golden Heights building, Chole Road",
		RegionCounty:  "Masaki, Dar es Salaam",
		CountryId:     defaultCountry.ID,
		ContactNumber: "+255 756 339 626",
		EmailAddress:  "info@afropavoanalytics.com",
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	db.Model(&companyOne).Omit("last_updated_by_id", "last_updated_at").Create(&companyOne)

	if logoContents != nil {
		db.Model(&model.CompanyFile{}).Create(&model.CompanyFile{
			ID:   companyOne.ID,
			Logo: logoContents,
		})
	}
	companyTwo := model.Company{
		Name:          "HSBC TA",
		CompanyTypeID: bankCompanyType.ID,
		Address:       "2nd Floor Golden Heights building, Chole Road",
		RegionCounty:  "Masaki, Dar es Salaam",
		CountryId:     defaultCountry.ID,
		ContactNumber: "+255 756 339 626",
		EmailAddress:  "info@afropavoanalytics.com",
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	db.Model(&companyTwo).Omit("last_updated_by_id", "last_updated_at").Create(&companyTwo)
	companyThree := model.Company{
		Name:          "AXIS BANK",
		CompanyTypeID: bankCompanyType.ID,
		Address:       "2nd Floor Golden Heights building, Chole Road",
		RegionCounty:  "Masaki, Dar es Salaam",
		CountryId:     defaultCountry.ID,
		ContactNumber: "+255 756 339 879",
		EmailAddress:  "info@axisbank.com",
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	db.Model(&companyThree).Omit("last_updated_by_id", "last_updated_at").Create(&companyThree)

	companyFour := model.Company{
		Name:          "FINANCE BANK",
		CompanyTypeID: bankCompanyType.ID,
		Address:       "2nd Floor Golden Heights building, Chole Road",
		RegionCounty:  "Masaki, Dar es Salaam",
		CountryId:     defaultCountry.ID,
		ContactNumber: "+255 756 339 879",
		EmailAddress:  "info@bank.com",
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	db.Model(&companyFour).Omit("last_updated_by_id", "last_updated_at").Create(&companyFour)

	customer_information_item := model.CustomerInformationItem{
		CustomerInformationID: 1,
		ModelVariableID:       2,
		Name:                  "Income",
		Value:                 50000,
		PreprocessedValue:     2}
	db.Model(&customer_information_item).Create(&customer_information_item)

	customer_information := model.CustomerInformation{

		ModelID:                        1,
		FirstName:                      "Sivansh",
		LastName:                       "Sai",
		CustomerIdProofNumber:          "8977463543854",
		CustomerIDProofTypeID:          1,
		CustomerIDProofType:            "National Identification Number",
		ContactNumber:                  "8977689666",
		City:                           "Hyderabad",
		ProbabilityOfDefaultPercentage: 87,
		GroupScore:                     178,
		CreatedBy:                      "Venkatraman",
		CreatedById:                    1,
		CompanyID:                      2,
		CreatedAt:                      timestamp,
	}

	db.Model(&customer_information).Create(&customer_information)
	customer_information1 := model.CustomerInformation{

		ModelID:                        1,
		FirstName:                      "Rohan",
		LastName:                       "Sai",
		CustomerIdProofNumber:          "857463543854",
		CustomerIDProofTypeID:          2,
		CustomerIDProofType:            "Passport",
		ContactNumber:                  "8977689666",
		City:                           "Hyderabad",
		ProbabilityOfDefaultPercentage: 87,
		GroupScore:                     178,
		CreatedBy:                      "Venkatraman",
		CreatedById:                    1,
		CompanyID:                      2,
		CreatedAt:                      timestamp,
	}

	db.Model(&customer_information1).Create(&customer_information1)
	db.Model(&customer_information).Create(&customer_information)
	customer_information4 := model.CustomerInformation{

		ModelID:                        4,
		FirstName:                      "Rohan",
		LastName:                       "Sai",
		CustomerIdProofNumber:          "857463543854",
		CustomerIDProofTypeID:          2,
		CustomerIDProofType:            "Passport",
		ContactNumber:                  "8977689666",
		City:                           "Hyderabad",
		ProbabilityOfDefaultPercentage: 87,
		GroupScore:                     178,
		CreatedBy:                      "Venkatraman",
		CreatedById:                    1,
		CompanyID:                      4,
		CreatedAt:                      timestamp,
	}

	db.Model(&customer_information4).Create(&customer_information4)
	customer_information2 := model.CustomerInformation{

		ModelID:                        3,
		FirstName:                      "Rehan",
		LastName:                       "Sai",
		CustomerIdProofNumber:          "857463543854",
		CustomerIDProofTypeID:          2,
		CustomerIDProofType:            "Passport",
		ContactNumber:                  "8977689666",
		City:                           "Chennai",
		ProbabilityOfDefaultPercentage: 60,
		GroupScore:                     148,
		CreatedBy:                      "Venkatraman",
		CreatedById:                    1,
		CompanyID:                      3,
		CreatedAt:                      timestamp,
	}

	db.Model(&customer_information2).Create(&customer_information2)

	customer_information3 := model.CustomerInformation{

		ModelID:                        2,
		FirstName:                      "Rehan",
		LastName:                       "Gurjar",
		CustomerIdProofNumber:          "857463543854",
		CustomerIDProofTypeID:          2,
		CustomerIDProofType:            "Passport",
		ContactNumber:                  "8977689666",
		City:                           "Chennai",
		ProbabilityOfDefaultPercentage: 80,
		GroupScore:                     148,
		CreatedBy:                      "Venkatraman",
		CreatedById:                    1,
		CompanyID:                      2,
		CreatedAt:                      timestamp,
	}

	db.Model(&customer_information3).Create(&customer_information3)

	modelOne := model.Model{
		CompanyID:      companyTwo.ID,
		Name:           "Model 1",
		Description:    "Model One",
		InterceptValue: -1.1743,
		Active:         true,
		CreatedById:    1,
		CreatedAt:      timestamp}
	db.Model(&modelOne).Omit("last_updated_by_id", "last_updated_at").Create(&modelOne)
	modelTwo := model.Model{
		CompanyID:      companyTwo.ID,
		Name:           "Model 2",
		Description:    "Model Two",
		InterceptValue: -1.1567,
		Active:         true,
		CreatedById:    1,
		CreatedAt:      timestamp}
	db.Model(&modelTwo).Omit("last_updated_by_id", "last_updated_at").Create(&modelTwo)
	modelThree := model.Model{
		CompanyID:      companyThree.ID,
		Name:           "Model 3",
		Description:    "Model Two",
		InterceptValue: -1.2345,
		Active:         true,
		CreatedById:    1,
		CreatedAt:      timestamp}
	db.Model(&modelThree).Omit("last_updated_by_id", "last_updated_at").Create(&modelThree)
	modelFour := model.Model{
		CompanyID:      companyFour.ID,
		Name:           "Model 4",
		Description:    "Model Two",
		InterceptValue: -1.5641,
		Active:         true,
		CreatedById:    1,
		CreatedAt:      timestamp}
	db.Model(&modelFour).Omit("last_updated_by_id", "last_updated_at").Create(&modelFour)

	mvEnumeration := &model.ModelVariable{
		ModelID:          modelOne.ID,
		Name:             "Gender",
		DataType:         2,
		CoefficientValue: 1}

	mvList := []*model.ModelVariable{
		{
			ModelID:                modelOne.ID,
			Name:                   "Group Allocatable Income",
			DataType:               1,
			CoefficientValue:       0.4617,
			MeanValue:              10807.92,
			StandardDeviationValue: 6250.856},
		{
			ModelID:                modelOne.ID,
			Name:                   "Income",
			DataType:               1,
			CoefficientValue:       0.6852,
			MeanValue:              3364.07,
			StandardDeviationValue: 1371.863},
		{
			ModelID:          modelOne.ID,
			Name:             "Income Ratio",
			DataType:         1,
			CoefficientValue: -0.7473},
		{
			ModelID:          modelOne.ID,
			Name:             "No of dependents",
			DataType:         1,
			CoefficientValue: 1.1567},
		{
			ModelID:          modelOne.ID,
			Name:             "Marital Status",
			DataType:         4,
			CoefficientValue: 1.1567},
		mvEnumeration}

	db.CreateInBatches(&mvList, 100)
	mvEnumeration1 := &model.ModelVariable{
		ModelID:          modelTwo.ID,
		Name:             "Gender",
		DataType:         2,
		CoefficientValue: 1}

	mvList1 := []*model.ModelVariable{
		{
			ModelID:                modelTwo.ID,
			Name:                   "Group Allocatable Income",
			DataType:               1,
			CoefficientValue:       0.4617,
			MeanValue:              10807.92,
			StandardDeviationValue: 6250.856},
		{
			ModelID:                modelTwo.ID,
			Name:                   "Income",
			DataType:               1,
			CoefficientValue:       0.6852,
			MeanValue:              3364.07,
			StandardDeviationValue: 1371.863},
		{
			ModelID:          modelTwo.ID,
			Name:             "Income Ratio",
			DataType:         1,
			CoefficientValue: -0.7473},
		{
			ModelID:          modelTwo.ID,
			Name:             "No of dependents",
			DataType:         1,
			CoefficientValue: 1.1567},
		{
			ModelID:          modelTwo.ID,
			Name:             "Marital Status",
			DataType:         4,
			CoefficientValue: 1.1567},
		mvEnumeration1}

	db.CreateInBatches(&mvList1, 100)
	mvEnumeration2 := &model.ModelVariable{
		ModelID:          modelThree.ID,
		Name:             "Gender",
		DataType:         2,
		CoefficientValue: 1}

	mvList2 := []*model.ModelVariable{
		{
			ModelID:                modelThree.ID,
			Name:                   "Group Allocatable Income",
			DataType:               1,
			CoefficientValue:       0.4617,
			MeanValue:              10807.92,
			StandardDeviationValue: 6250.856},
		{
			ModelID:                modelThree.ID,
			Name:                   "Income",
			DataType:               1,
			CoefficientValue:       0.6852,
			MeanValue:              3364.07,
			StandardDeviationValue: 1371.863},
		{
			ModelID:          modelThree.ID,
			Name:             "Income Ratio",
			DataType:         1,
			CoefficientValue: -0.7473},
		{
			ModelID:          modelThree.ID,
			Name:             "No of dependents",
			DataType:         1,
			CoefficientValue: 1.1567},
		{
			ModelID:          modelThree.ID,
			Name:             "Marital Status",
			DataType:         4,
			CoefficientValue: 1.1567},
		mvEnumeration2}

	db.CreateInBatches(&mvList2, 100)
	mvEnumeration3 := &model.ModelVariable{
		ModelID:          modelFour.ID,
		Name:             "Gender",
		DataType:         2,
		CoefficientValue: 1}

	mvList3 := []*model.ModelVariable{
		{
			ModelID:                modelFour.ID,
			Name:                   "Group Allocatable Income",
			DataType:               1,
			CoefficientValue:       0.4617,
			MeanValue:              10807.92,
			StandardDeviationValue: 6250.856},
		{
			ModelID:                modelFour.ID,
			Name:                   "Income",
			DataType:               1,
			CoefficientValue:       0.6852,
			MeanValue:              3364.07,
			StandardDeviationValue: 1371.863},
		{
			ModelID:          modelFour.ID,
			Name:             "Income Ratio",
			DataType:         1,
			CoefficientValue: -0.7473},
		{
			ModelID:          modelFour.ID,
			Name:             "No of dependents",
			DataType:         1,
			CoefficientValue: 1.1567},
		{
			ModelID:          modelFour.ID,
			Name:             "Marital Status",
			DataType:         4,
			CoefficientValue: 1.1567},
		mvEnumeration3}

	db.CreateInBatches(&mvList3, 100)
	db.CreateInBatches([]model.EnumerationItem{
		{
			ModelVariableID: mvEnumeration.ID,
			Text:            "Male",
			Value:           0},
		{
			ModelVariableID: mvEnumeration.ID,
			Text:            "Female",
			Value:           1}}, 100)

	urList = []model.UserRole{{Name: "CRA"}}
	db.CreateInBatches(&urList, 100)

	dob := datatypes.Date{}
	t, _ := time.Parse(DtFormatSimpleDate, "1976-08-15")
	dob.Scan(t)

	user := model.User{
		FirstName:     "Derick",
		LastName:      "Kazimoto",
		Email:         "d.kazimoto@afropavoanalytics.com",
		Password:      "hello123",
		ContactNumber: "+255 756 339 626",
		CompanyId:     companyOne.ID,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Venkatraman",
		LastName:      "Balasubramanian",
		Email:         "venkatraman@phiaura.in",
		Password:      "hello123",
		ContactNumber: "+91 96000 02859",
		CompanyId:     1,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Rwebugisa",
		LastName:      "Mutahaba",
		Email:         "r.mutahaba@afropavoanalytics.com",
		Password:      "tausiadmin123",
		ContactNumber: "+255 765 664 084",
		CompanyId:     1,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Sophia",
		LastName:      "Mwinyi",
		Email:         "s.mwinyi@afropavoanalytics.com",
		Password:      "tausiadmin123",
		ContactNumber: "+255 713 557 731",
		CompanyId:     1,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Prakash - QA",
		LastName:      "d",
		Email:         "prakash.d@3edge.in",
		Password:      "hello123",
		ContactNumber: "+91 96000 02859",
		CompanyId:     1,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Muralikrishnan",
		LastName:      "s",
		Email:         "muralikrishnan.s@3edge.in",
		Password:      "test123",
		ContactNumber: "+91 9500984141",
		CompanyId:     2,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)

	user = model.User{
		FirstName:     "Muralikrishnan",
		LastName:      "s",
		Email:         "ssmkrishnan86@gmail.com",
		Password:      "test123",
		ContactNumber: "+91 9500984141",
		CompanyId:     2,
		UserRoles:     urList,
		EmailVerified: true,
		Active:        true,
		CreatedById:   1,
		CreatedAt:     timestamp}
	create(db, &user)
}

func loadCountries(defaultCountry *model.Country, db *gorm.DB) (err error) {
	var cList = []*model.Country{{Name: "Åland Islands"},
		{Name: "Albania"},
		{Name: "Algeria"},
		{Name: "American Samoa"},
		{Name: "Andorra"},
		{Name: "Angola"},
		{Name: "Anguilla"},
		{Name: "Antarctica"},
		{Name: "Antigua and Barbuda"},
		{Name: "Argentina"},
		{Name: "Armenia"},
		{Name: "Aruba"},
		{Name: "Australia"},
		{Name: "Austria"},
		{Name: "Azerbaijan"},
		{Name: "Bahamas"},
		{Name: "Bahrain"},
		{Name: "Bangladesh"},
		{Name: "Barbados"},
		{Name: "Belarus"},
		{Name: "Belgium"},
		{Name: "Belize"},
		{Name: "Benin"},
		{Name: "Bermuda"},
		{Name: "Bhutan"},
		{Name: "Bolivia (Plurinational State of)"},
		{Name: "Bonaire, Sint Eustatius and Saba"},
		{Name: "Bosnia and Herzegovina"},
		{Name: "Botswana"},
		{Name: "Bouvet Island"},
		{Name: "Brazil"},
		{Name: "British Indian Ocean Territory"},
		{Name: "Brunei Darussalam"},
		{Name: "Bulgaria"},
		{Name: "Burkina Faso"},
		{Name: "Burundi"},
		{Name: "Cabo Verde"},
		{Name: "Cambodia"},
		{Name: "Cameroon"},
		{Name: "Canada"},
		{Name: "Cayman Islands"},
		{Name: "Central African Republic"},
		{Name: "Chad"},
		{Name: "Chile"},
		{Name: "China"},
		{Name: "Christmas Island"},
		{Name: "Cocos (Keeling) Islands"},
		{Name: "Colombia"},
		{Name: "Comoros"},
		{Name: "Congo"},
		{Name: "Congo, Democratic Republic of the"},
		{Name: "Cook Islands"},
		{Name: "Costa Rica"},
		{Name: "Côte d'Ivoire"},
		{Name: "Croatia"},
		{Name: "Cuba"},
		{Name: "Curaçao"},
		{Name: "Cyprus"},
		{Name: "Czechia"},
		{Name: "Denmark"},
		{Name: "Djibouti"},
		{Name: "Dominica"},
		{Name: "Dominican Republic"},
		{Name: "Ecuador"},
		{Name: "Egypt"},
		{Name: "El Salvador"},
		{Name: "Equatorial Guinea"},
		{Name: "Eritrea"},
		{Name: "Estonia"},
		{Name: "Eswatini"},
		{Name: "Ethiopia"},
		{Name: "Falkland Islands (Malvinas)"},
		{Name: "Faroe Islands"},
		{Name: "Fiji"},
		{Name: "Finland"},
		{Name: "France"},
		{Name: "French Guiana"},
		{Name: "French Polynesia"},
		{Name: "French Southern Territories"},
		{Name: "Gabon"},
		{Name: "Gambia"},
		{Name: "Georgia"},
		{Name: "Germany"},
		{Name: "Ghana"},
		{Name: "Gibraltar"},
		{Name: "Greece"},
		{Name: "Greenland"},
		{Name: "Grenada"},
		{Name: "Guadeloupe"},
		{Name: "Guam"},
		{Name: "Guatemala"},
		{Name: "Guernsey"},
		{Name: "Guinea"},
		{Name: "Guinea-Bissau"},
		{Name: "Guyana"},
		{Name: "Haiti"},
		{Name: "Heard Island and McDonald Islands"},
		{Name: "Holy See"},
		{Name: "Honduras"},
		{Name: "Hong Kong"},
		{Name: "Hungary"},
		{Name: "Iceland"},
		{Name: "India"},
		{Name: "Indonesia"},
		{Name: "Iran (Islamic Republic of)"},
		{Name: "Iraq"},
		{Name: "Ireland"},
		{Name: "Isle of Man"},
		{Name: "Israel"},
		{Name: "Italy"},
		{Name: "Jamaica"},
		{Name: "Japan"},
		{Name: "Jersey"},
		{Name: "Jordan"},
		{Name: "Kazakhstan"},
		{Name: "Kenya"},
		{Name: "Kiribati"},
		{Name: "Korea (Democratic People's Republic of)"},
		{Name: "Korea, Republic of"},
		{Name: "Kuwait"},
		{Name: "Kyrgyzstan"},
		{Name: "Lao People's Democratic Republic"},
		{Name: "Latvia"},
		{Name: "Lebanon"},
		{Name: "Lesotho"},
		{Name: "Liberia"},
		{Name: "Libya"},
		{Name: "Liechtenstein"},
		{Name: "Lithuania"},
		{Name: "Luxembourg"},
		{Name: "Macao"},
		{Name: "Madagascar"},
		{Name: "Malawi"},
		{Name: "Malaysia"},
		{Name: "Maldives"},
		{Name: "Mali"},
		{Name: "Malta"},
		{Name: "Marshall Islands"},
		{Name: "Martinique"},
		{Name: "Mauritania"},
		{Name: "Mauritius"},
		{Name: "Mayotte"},
		{Name: "Mexico"},
		{Name: "Micronesia (Federated States of)"},
		{Name: "Moldova, Republic of"},
		{Name: "Monaco"},
		{Name: "Mongolia"},
		{Name: "Montenegro"},
		{Name: "Montserrat"},
		{Name: "Morocco"},
		{Name: "Mozambique"},
		{Name: "Myanmar"},
		{Name: "Namibia"},
		{Name: "Nauru"},
		{Name: "Nepal"},
		{Name: "Netherlands"},
		{Name: "New Caledonia"},
		{Name: "New Zealand"},
		{Name: "Nicaragua"},
		{Name: "Niger"},
		{Name: "Nigeria"},
		{Name: "Niue"},
		{Name: "Norfolk Island"},
		{Name: "North Macedonia"},
		{Name: "Northern Mariana Islands"},
		{Name: "Norway"},
		{Name: "Oman"},
		{Name: "Pakistan"},
		{Name: "Palau"},
		{Name: "Palestine, State of"},
		{Name: "Panama"},
		{Name: "Papua New Guinea"},
		{Name: "Paraguay"},
		{Name: "Peru"},
		{Name: "Philippines"},
		{Name: "Pitcairn"},
		{Name: "Poland"},
		{Name: "Portugal"},
		{Name: "Puerto Rico"},
		{Name: "Qatar"},
		{Name: "Réunion"},
		{Name: "Romania"},
		{Name: "Russian Federation"},
		{Name: "Rwanda"},
		{Name: "Saint Barthélemy"},
		{Name: "Saint Helena, Ascension and Tristan da Cunha"},
		{Name: "Saint Kitts and Nevis"},
		{Name: "Saint Lucia"},
		{Name: "Saint Martin (French part)"},
		{Name: "Saint Pierre and Miquelon"},
		{Name: "Saint Vincent and the Grenadines"},
		{Name: "Samoa"},
		{Name: "San Marino"},
		{Name: "Sao Tome and Principe"},
		{Name: "Saudi Arabia"},
		{Name: "Senegal"},
		{Name: "Serbia"},
		{Name: "Seychelles"},
		{Name: "Sierra Leone"},
		{Name: "Singapore"},
		{Name: "Sint Maarten (Dutch part)"},
		{Name: "Slovakia"},
		{Name: "Slovenia"},
		{Name: "Solomon Islands"},
		{Name: "Somalia"},
		{Name: "South Africa"},
		{Name: "South Georgia and the South Sandwich Islands"},
		{Name: "South Sudan"},
		{Name: "Spain"},
		{Name: "Sri Lanka"},
		{Name: "Sudan"},
		{Name: "Suriname"},
		{Name: "Svalbard and Jan Mayen"},
		{Name: "Sweden"},
		{Name: "Switzerland"},
		{Name: "Syrian Arab Republic"},
		{Name: "Taiwan, Province of China"},
		{Name: "Tajikistan"},
		defaultCountry,
		{Name: "Thailand"},
		{Name: "Timor-Leste"},
		{Name: "Togo"},
		{Name: "Tokelau"},
		{Name: "Tonga"},
		{Name: "Trinidad and Tobago"},
		{Name: "Tunisia"},
		{Name: "Turkey"},
		{Name: "Turkmenistan"},
		{Name: "Turks and Caicos Islands"},
		{Name: "Tuvalu"},
		{Name: "Uganda"},
		{Name: "Ukraine"},
		{Name: "United Arab Emirates"},
		{Name: "United Kingdom of Great Britain and Northern Ireland"},
		{Name: "United States of America"},
		{Name: "United States Minor Outlying Islands"},
		{Name: "Uruguay"},
		{Name: "Uzbekistan"},
		{Name: "Vanuatu"},
		{Name: "Venezuela (Bolivarian Republic of)"},
		{Name: "Viet Nam"},
		{Name: "Virgin Islands (British)"},
		{Name: "Virgin Islands (U.S.)"},
		{Name: "Wallis and Futuna"},
		{Name: "Western Sahara"},
		{Name: "Yemen"},
		{Name: "Zambia"},
		{Name: "Zimbabwe"}}

	db.CreateInBatches(&cList, 100)
	err = db.Error

	return
}

func loadCities(db *gorm.DB) (err error) {
	var cList = []*model.City{{Name: "Dar es Salaam"},
		{Name: "Mwanza"},
		{Name: "Zanzibar City"},
		{Name: "Arusha"},
		{Name: "Mbeya"},
		{Name: "Morogoro"},
		{Name: "Tanga"},
		{Name: "Dodoma"},
		{Name: "Kigoma"},
		{Name: "Moshi"},
		{Name: "Tabora"},
		{Name: "Songea"},
		{Name: "Musoma"},
		{Name: "Shinyanga"},
		{Name: "Katumba"},
		{Name: "Iringa"},
		{Name: "Ushirombo"},
		{Name: "Mtwara"},
		{Name: "Kilosa"},
		{Name: "Sumbawanga"},
		{Name: "Bagamoyo"},
		{Name: "Mpanda"},
		{Name: "Bukoba"},
		{Name: "Singida Urban"},
		{Name: "Uyovu"},
		{Name: "Sengerama"},
		{Name: "Kalangalala"},
		{Name: "Mishoma"},
		{Name: "Mererani"},
		{Name: "Buseresere"},
		{Name: "Bunda"},
		{Name: "Makambako"},
		{Name: "Katoro"},
		{Name: "Ifakara"},
		{Name: "Njombe"},
		{Name: "Utengule Usongwe"},
		{Name: "Kiranyi"},
		{Name: "Siha Kati"},
		{Name: "Nkome"},
		{Name: "Nkololo"},
		{Name: "Nguruka"},
		{Name: "Lindi"},
		{Name: "Vwawa"}}
	db.CreateInBatches(&cList, 100)
	err = db.Error
	return
}

func create(db *gorm.DB, user *model.User) {
	passwordText := user.Password
	user.Password = "123"

	db.Model(&user).Omit("last_updated_by_id", "last_updated_at").Create(&user)
	user.Password = passwordText

	err := user.Hash()
	if err != nil {
		panic(err)
	}

	err = db.Model(&user).UpdateColumn("password", user.Password).Error
	if err != nil {
		panic(err)
	}
	err = db.Model(&user).Association("UserRoles").Clear()
	if err != nil {
		panic(err)
	}
}
