package company

// CompanyType 公司類型
type CompanyType string

const (
	Startup       CompanyType = "startup"        // 新創公司
	Traditional   CompanyType = "traditional"    // 傳產公司
	CloudReseller CompanyType = "cloud_reseller" // 雲端代理商
	CloudNative   CompanyType = "cloud_native"   // 雲端公司
)

// Company 公司實體
type Company struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Type            CompanyType `json:"type"`
	Capital         int64       `json:"capital"`          // 資本額 (萬元)
	Employees       int         `json:"employees"`        // 員工人數
	IsInternational bool        `json:"is_international"` // 是否跨國企業
	ProductCycle    string      `json:"product_cycle"`    // 產品週期: "development", "launch", "growth", "mature"
	TechDebt        int         `json:"tech_debt"`        // 技術債 (影響決策)
	SecurityLevel   int         `json:"security_level"`   // 資安等級 1-5
	CloudAdoption   float64     `json:"cloud_adoption"`   // 雲端採用率 0-100%
	Infrastructure  []string    `json:"infrastructure"`   // 已部署的基礎設施
}

// CompanyDefault 公司預設屬性
type CompanyDefault struct {
	Capital       int64
	Employees     int
	SecurityLevel int
	CloudAdoption float64
	ProductCycle  string
	TechDebt      int
}
