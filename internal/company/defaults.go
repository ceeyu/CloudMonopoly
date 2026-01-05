package company

// CompanyDefaults 各公司類型初始屬性
var CompanyDefaults = map[CompanyType]CompanyDefault{
	Startup: {
		Capital:       500,
		Employees:     10,
		SecurityLevel: 2,
		CloudAdoption: 30,
		ProductCycle:  "development",
		TechDebt:      20,
	},
	Traditional: {
		Capital:       5000,
		Employees:     200,
		SecurityLevel: 3,
		CloudAdoption: 10,
		ProductCycle:  "mature",
		TechDebt:      60,
	},
	CloudReseller: {
		Capital:       2000,
		Employees:     50,
		SecurityLevel: 4,
		CloudAdoption: 80,
		ProductCycle:  "growth",
		TechDebt:      30,
	},
	CloudNative: {
		Capital:       1000,
		Employees:     30,
		SecurityLevel: 4,
		CloudAdoption: 95,
		ProductCycle:  "launch",
		TechDebt:      10,
	},
}

// ValidCompanyTypes 有效的公司類型列表
var ValidCompanyTypes = []CompanyType{
	Startup,
	Traditional,
	CloudReseller,
	CloudNative,
}
