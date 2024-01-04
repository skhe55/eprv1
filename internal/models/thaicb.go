package models

type ThaiCBResponse struct {
	Result struct {
		Timestamp string `json:"timestamp"`
		API       string `json:"api"`
		Data      struct {
			DataHeader struct {
				ReportNameEng      string `json:"report_name_eng"`
				ReportNameTh       string `json:"report_name_th"`
				ReportUoqNameEng   string `json:"report_uoq_name_eng"`
				ReportUoqNameTh    string `json:"report_uoq_name_th"`
				ReportSourceOfData []struct {
					SourceOfDataEng string `json:"source_of_data_eng"`
					SourceOfDataTh  string `json:"source_of_data_th"`
				} `json:"report_source_of_data"`
				ReportRemark []struct {
					ReportRemarkEng string `json:"report_remark_eng"`
					ReportRemarkTh  string `json:"report_remark_th"`
				} `json:"report_remark"`
				LastUpdated string `json:"last_updated"`
			} `json:"data_header"`
			DataDetail []struct {
				Period          string `json:"period"`
				CurrencyID      string `json:"currency_id"`
				CurrencyNameTh  string `json:"currency_name_th"`
				CurrencyNameEng string `json:"currency_name_eng"`
				BuyingSight     string `json:"buying_sight"`
				BuyingTransfer  string `json:"buying_transfer"`
				Selling         string `json:"selling"`
				MidRate         string `json:"mid_rate"`
			} `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}
