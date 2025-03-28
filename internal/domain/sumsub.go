package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SumsubDataInfoAddress struct {
	SubStreet        string `json:"subStreet"`
	SubStreetEn      string `json:"subStreetEn"`
	Street           string `json:"street"`
	StreetEn         string `json:"streetEn"`
	State            string `json:"state"`
	StateEn          string `json:"stateEn"`
	Town             string `json:"town"`
	TownEn           string `json:"townEn"`
	PostCode         string `json:"postCode"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formattedAddress"`
	// LocationPosition interface{} `json:"locationPosition"`
}

type SumsubDataInfoIdDocsAddress struct {
	Street           string `json:"street"`
	StreetEn         string `json:"streetEn"`
	Town             string `json:"town"`
	TownEn           string `json:"townEn"`
	Country          string `json:"country"`
	FormattedAddress string `json:"formattedAddress"`
	// LocationPosition interface{} `json:"locationPosition"`
}

type SumsubDataInfoIdDocs struct {
	IdDocType    string                      `json:"idDocType"`
	Country      string                      `json:"country"`
	FirstName    string                      `json:"firstName"`
	FirstNameEn  string                      `json:"firstNameEn"`
	LastName     string                      `json:"lastName"`
	LastNameEn   string                      `json:"lastNameEn"`
	IssuedDate   string                      `json:"issuedDate"`
	Number       string                      `json:"number"`
	Dob          string                      `json:"dob"`
	Gender       string                      `json:"gender"`
	PlaceOfBirth string                      `json:"placeOfBirth"`
	Address      SumsubDataInfoIdDocsAddress `json:"address"`
}

type SumsubDataInfo struct {
	FirstName      string                  `json:"firstName"`
	FirstNameEn    string                  `json:"firstNameEn"`
	LastName       string                  `json:"lastName"`
	LastNameEn     string                  `json:"lastNameEn"`
	Dob            string                  `json:"dob"`
	Gender         string                  `json:"gender"`
	PlaceOfBirth   string                  `json:"placeOfBirth"`
	PlaceOfBirthEn string                  `json:"placeOfBirthEn"`
	Country        string                  `json:"country"`
	Addresses      []SumsubDataInfoAddress `json:"addresses"`
	IdDocs         []SumsubDataInfoIdDocs  `json:"idDocs"`
}

type SumsubDataAgreement struct {
	CreatedAt  string   `json:"createdAt"`
	AcceptedAt string   `json:"acceptedAt"`
	Source     string   `json:"source"`
	RecordIds  []string `json:"recordIds"`
}

type SumsubDataRequiredIdDocsDocSets struct {
	IdDocSetType  string   `json:"idDocSetType"`
	Types         []string `json:"types"`
	VideoRequired string   `json:"videoRequired"`
	CaptureMode   string   `json:"captureMode"`
	UploaderMode  string   `json:"uploaderMode"`
}

type SumsubDataRequiredIdDocs struct {
	IncludedCountries []string                          `json:"includedCountries"`
	DocSets           []SumsubDataRequiredIdDocsDocSets `json:"docSets"`
}

type SumsubDataReviewReviewResult struct {
	ReviewAnswer string `json:"reviewAnswer"`
}

type SumsubDataReview struct {
	ReviewId              string `json:"reviewId"`
	AttemptId             string `json:"attemptId"`
	AttemptCnt            int    `json:"attemptCnt"`
	ElapsedSincePendingMs int    `json:"elapsedSincePendingMs"`
	ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs"`
	Reprocessing          bool   `json:"reprocessing"`
	LevelName             string `json:"levelName"`
	// LevelAutoCheckMode    interface{}                  `json:"levelAutoCheckMode"`
	CreateDate   string                       `json:"createDate"`
	ReviewDate   string                       `json:"reviewDate"`
	ReviewResult SumsubDataReviewReviewResult `json:"reviewResult"`
	ReviewStatus string                       `json:"reviewStatus"`
	Priority     int                          `json:"priority"`
}

type SumsubDataRiskLabels struct {
	AttemptId string   `json:"attemptId"`
	CreatedAt string   `json:"createdAt"`
	Device    []string `json:"device"`
	Selfie    []string `json:"selfie"`
}

type SumsubData struct {
	Id                string                     `json:"id"`
	CreatedAt         string                     `json:"createdAt"`
	CreatedBy         string                     `json:"createdBy"`
	Key               string                     `json:"key"`
	ClientId          string                     `json:"clientId"`
	InspectionId      string                     `json:"inspectionId"`
	ExternalUserId    string                     `json:"externalUserId"`
	Info              SumsubDataInfo             `json:"info"`
	ApplicantPlatform string                     `json:"applicantPlatform"`
	IpCountry         string                     `json:"ipCountry"`
	AuthCode          string                     `json:"authCode"`
	Agreement         SumsubDataAgreement        `json:"agreement"`
	RequiredIdDocs    SumsubDataRequiredIdDocs   `json:"requiredIdDocs"`
	Review            SumsubDataReview           `json:"review"`
	Lang              string                     `json:"lang"`
	Type              string                     `json:"type"`
	RiskLabels        SumsubDataRiskLabels       `json:"riskLabels"`
	Email             string                     `json:"email"`
	Phone             string                     `json:"phone"`
	Questionnaires    []SumsubDataQuestionnaires `json:"questionnaires"`
}

type SumsubDataQuestionnaires struct {
	Id       string                           `json:"id"`
	Sections SumsubDataQuestionnairesSections `json:"sections"`
}

type SumsubDataQuestionnairesSections struct {
	JiBenZiXun SumsubDataQuestionnairesSectionsJiBenZiXun `json:"jiBenZiXun"`
}

type SumsubDataQuestionnairesSectionsJiBenZiXun struct {
	Items SumsubDataQuestionnairesSectionsJiBenZiXunItem `json:"items"`
}

type SumsubDataQuestionnairesSectionsJiBenZiXunItem struct {
	NinDeZhiYe SumsubDataQuestionnairesSectionsJiBenZiXunItemNinDeZhiYe `json:"ninDeZhiYe"`
}

type SumsubDataQuestionnairesSectionsJiBenZiXunItemNinDeZhiYe struct {
	Value string `json:"value"`
}

type SumsubDataResApplicantIdDoc struct {
	IdDocType    string `json:"idDocType,omitempty"`
	Country      string `json:"country,omitempty"`
	FirstName    string `json:"firstName,omitempty"`
	FirstNameEn  string `json:"firstNameEn,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	MiddleNameEn string `json:"middleNameEn,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	LastNameEn   string `json:"lastNameEn,omitempty"`
	DateOfBirth  string `json:"dob,omitempty"` // yyyy-mm-dd format
	Number       string `json:"number,omitempty"`
	IssuedDate   string `json:"issuedDate,omitempty"` // yyyy-mm-dd format
}

type SumsubDataResApplicantInfo struct {
	FirstName    string                        `json:"firstName,omitempty"`
	FirstNameEn  string                        `json:"firstNameEn,omitempty"`
	MiddleName   string                        `json:"middleName,omitempty"`
	MiddleNameEn string                        `json:"middleNameEn,omitempty"`
	LastName     string                        `json:"lastName,omitempty"`
	LastNameEn   string                        `json:"lastNameEn,omitempty"`
	Dob          string                        `json:"dob,omitempty"` // yyyy-mm-dd format
	Gender       string                        `json:"gender,omitempty"`
	Country      string                        `json:"country,omitempty"`
	IdDocs       []SumsubDataResApplicantIdDoc `json:"idDocs,omitempty"`
}

type SumsubDataResApplicant struct {
	ID             string                     `json:"id,omitempty"`
	CreatedAt      string                     `json:"createdAt,omitempty"`
	Key            string                     `json:"key,omitempty"`
	ClientID       string                     `json:"clientId,omitempty"`
	InspectionID   string                     `json:"inspectionId,omitempty"`
	ExternalUserID string                     `json:"externalUserId,omitempty"`
	Info           SumsubDataResApplicantInfo `json:"info,omitempty"`
	FixedInfo      SumsubDataResApplicantInfo `json:"fixedInfo,omitempty"`
	Phone          string                     `json:"phone,omitempty"`
	Review         struct {
		ElapsedSincePendingMs int    `json:"elapsedSincePendingMs,omitempty"`
		ElapsedSinceQueuedMs  int    `json:"elapsedSinceQueuedMs,omitempty"`
		Reprocessing          bool   `json:"reprocessing,omitempty"`
		CreateDate            string `json:"createDate,omitempty"`
		ReviewDate            string `json:"reviewDate,omitempty"`
		StartDate             string `json:"startDate,omitempty"`
		ReviewResult          struct {
			ReviewAnswer      string   `json:"reviewAnswer,omitempty"`
			RejectLabels      []string `json:"rejectLabels,omitempty"`
			ModerationComment string   `json:"moderationComment,omitempty"`
			ClientComment     string   `json:"clientComment,omitempty"`
		} `json:"reviewResult,omitempty"`
		ReviewStatus           string `json:"reviewStatus,omitempty"`
		NotificationFailureCnt int    `json:"notificationFailureCnt,omitempty"`
		Priority               int    `json:"priority,omitempty"`
		LevelName              string `json:"levelName,omitempty"`
	} `json:"review,omitempty"`
	Lang       string   `json:"lang,omitempty"`
	Type       string   `json:"type,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	RiskLabels struct {
		Email      []string `json:"email,omitempty"`
		Phone      []string `json:"phone,omitempty"`
		Device     []string `json:"device,omitempty"`
		CrossCheck []string `json:"crossCheck,omitempty"`
		Selfie     []string `json:"selfie,omitempty"`
		Aml        []string `json:"aml,omitempty"`
		Person     []string `json:"person,omitempty"`
	} `json:"riskLabels,omitempty"`
	Questionnaires []struct {
		Id       string  `json:"id,omitempty"`
		Score    float64 `json:"score,omitempty"`
		Sections struct {
			BasicInfo struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					LastName struct {
						Value string `json:"value,omitempty"`
					} `json:"lastName,omitempty"`
					FirstName struct {
						Value string `json:"value,omitempty"`
					} `json:"firstName,omitempty"`
					IDCategory struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"ID_category,omitempty"`
					IDLocation struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"ID_location,omitempty"`
					Adress struct {
						Value string `json:"value,omitempty"`
					} `json:"adress,omitempty"`
				}
			} `json:"Basic_info,omitempty"`
			JiBenZiLiao struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					KaiHuMuDe struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"kaiHuMuDe,omitempty"`
					NianShouRuJiJu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"nianShouRuJiJu,omitempty"`
					ZhiCheng struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhiCheng,omitempty"`
					ZhiYe struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhiYe,omitempty"`
					XianRenZhiGongSiHuoJ struct {
						Value string `json:"value,omitempty"`
					} `json:"xianRenZhiGongSiHuoJ,omitempty"`
					RenZhiNianShu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"renZhiNianShu,omitempty"`
					ZhangHaoYongTu struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"zhangHaoYongTu,omitempty"`
					RenZhiGongSiLeiXing struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"renZhiGongSiLeiXing,omitempty"`
				}
			} `json:"jiBenZiLiao,omitempty"`
			GaoZhiShiXiang struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					ShuoMingNeiRong struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"shuoMingNeiRong,omitempty"`
				}
			} `json:"gaoZhiShiXiang,omitempty"`
			ZiChanPingGu struct {
				Score float64 `json:"score,omitempty"`
				Items struct {
					GeRenZongZiChan struct {
						Value string  `json:"value,omitempty"`
						Score float64 `json:"score,omitempty"`
					} `json:"geRenZongZiChan,omitempty"`
					ZiJinLaiYuan struct {
						Values []string `json:"values,omitempty"`
						Score  float64  `json:"score,omitempty"`
					} `json:"ziJinLaiYuan,omitempty"`
					ZiChanZhengMingWenJi struct {
						Values []string `json:"values,omitempty"`
					} `json:"ziChanZhengMingWenJi,omitempty"`
				}
			} `json:"ziChanPingGu,omitempty"`
		}
	} `json:"questionnaires,omitempty"`
}

type SumsubDataApplicantRequiredIdDocs struct {
	IDENTITY struct {
		ReviewResult struct {
			ModerationComment string `json:"moderationComment"`
			ReviewAnswer      string `json:"reviewAnswer"`
		} `json:"reviewResult"`
		Country   string `json:"country"`
		IdDocType string `json:"idDocType"`
		ImageIds  []int  `json:"imageIds"`
	} `json:"IDENTITY"`
	SELFIE struct {
		ReviewResult struct {
			ReviewAnswer string `json:"reviewAnswer"`
		} `json:"reviewResult"`
		Country   string `json:"country"`
		IdDocType string `json:"idDocType"`
		ImageIds  []int  `json:"imageIds"`
	} `json:"SELFIE"`
}

type SumsubImageDos struct {
	IdCardFrontImgId string
	IdCardBackImgId  string
	SelfieImgId      string
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s *SumsubData) Scan(value interface{}) error {
	if value == nil {
		*s = SumsubData{}
		return nil
	}

	bytesValue, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan SumsubData")
	}

	data := SumsubData{}
	if err := json.Unmarshal(bytesValue, &data); err != nil {
		return errors.New("invalid scan SumsubData unmarshal")
	}

	*s = data
	return nil
}

// Value return json value, implement driver.Valuer interface
func (s SumsubData) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type SumsubTag string

func (e SumsubTag) String() string { return string(e) }

const (
	SumsubTagCib SumsubTag = "告誡名單"
)

type SumsubWebhookType string

func (e SumsubWebhookType) String() string { return string(e) }

const (
	SumsubWebhookTypeApplicantReviewed    SumsubWebhookType = "applicantReviewed"
	SumsubWebhookTypeApplicantTagsChanged SumsubWebhookType = "applicantTagsChanged"
	SumsubWebhookTypeApplicantOnHold      SumsubWebhookType = "applicantOnHold"
)

type SumsubApplicantStatus string

func (e SumsubApplicantStatus) String() string { return string(e) }

const (
	SumsubApplicantStatusCompleted SumsubApplicantStatus = "completed"
	SumsubApplicantStatusOnHold    SumsubApplicantStatus = "onHold"
)

var KYCDocOccupation = map[string]string{
	"1":  "農牧漁林業",
	"2":  "礦業採石業",
	"3":  "交通運輸業",
	"4":  "餐旅業",
	"5":  "建築工程業",
	"6":  "製造業",
	"7":  "新聞廣告業",
	"8":  "醫療保健業",
	"9":  "娛樂業",
	"10": "文教業",
	"11": "宗教團體",
	"12": "公職人員",
	"13": "金融服務業",
	"14": "資訊業",
	"15": "半導體業",
	"16": "其他服務業",
	"17": "殯葬業",
	"18": "進出口貿易業",
	"19": "航空業",
	"20": "軍事人員",
	"21": "法律人士",
	"22": "家管",
	"23": "學生",
	"24": "退休人士",
	"25": "自由業",
	"26": "待業中",
	"27": "其他（未列表）",
}

type SumsubHistoryReviewData struct {
	ApplicantId       string              `json:"applicantId"`
	ApplicantActionId *string             `json:"applicantActionId"`
	Items             []SumsubHistoryItem `json:"items"`
}

type SumsubHistoryItem struct {
	Id                 string                        `json:"id"`
	Ts                 string                        `json:"ts"`
	Activity           string                        `json:"activity"`
	SubjectName        string                        `json:"subjectName"`
	Country            string                        `json:"country"`
	Status             string                        `json:"status"`
	ReviewAnswer       string                        `json:"reviewAnswer"`
	ReviewRejectType   string                        `json:"reviewRejectType"`
	ReviewResult       SumsubHistoryItemReviewResult `json:"reviewResult"`
	ApplicantLevelName string                        `json:"applicantLevelName"`
	AttemptId          string                        `json:"attemptId"`
	ApplicantId        string                        `json:"applicantId"`
	ClientId           string                        `json:"clientId"`
}

type SumsubHistoryItemReviewResult struct {
	RejectLabels []string `json:"rejectLabels"`
}

var KYCTag = map[string]string{
	"APPLICANT_INTERRUPTED_INTERVIEW":       "在視頻識別電話中，申請人拒絕結束面試。",
	"ADDITIONAL_DOCUMENT_REQUIRED":          "透過檢查所需的其他文件。",
	"BACK_SIDE_MISSING":                     "文件背面缺失。",
	"BAD_AVATAR":                            "頭像不符合客戶的要求。",
	"BAD_FACE_MATCHING":                     "文檔和自拍照之間的面部檢查失敗。",
	"BAD_PROOF_OF_ADDRESS":                  "申請人上傳了錯誤的地址證明。",
	"BAD_PROOF_OF_IDENTITY":                 "申請人上傳了錯誤的身份證件。",
	"BAD_PROOF_OF_PAYMENT":                  "申請人上傳了錯誤的付款證明。",
	"BAD_SELFIE":                            "申請人上傳了一張糟糕的自拍照。",
	"BAD_VIDEO_SELFIE":                      "申請人上傳了一張糟糕的自拍照視頻。",
	"BLACK_AND_WHITE":                       "申請人上傳了證件的黑白照片。",
	"COMPANY_NOT_DEFINED_BENEFICIARIES":     "無法識別並適當核實該實體的受益所有人。",
	"COMPANY_NOT_DEFINED_REPRESENTATIVES":   "代表未定義。",
	"COMPANY_NOT_DEFINED_STRUCTURE":         "無法建立實體控制結構。",
	"COMPANY_NOT_VALIDATED_BENEFICIARIES":   "受益人未經驗證。",
	"COMPANY_NOT_VALIDATED_REPRESENTATIVES": "代表未經過驗證。",
	"CONNECTION_INTERRUPTED":                "視頻識別呼叫連接已中斷。",
	"DIGITAL_DOCUMENT":                      "申請人上傳了該文件的電子版。",
	"DOCUMENT_DEPRIVED":                     "申請人已被剝奪該文件。",
	"DOCUMENT_DAMAGED":                      "文檔已損壞。",
	"DOCUMENT_MISSING":                      "在視頻識別電話中，申請人拒絕出示或沒有所需文件。",
	"DOCUMENT_PAGE_MISSING":                 "文檔的某些頁面丟失。",
	"EXPIRATION_DATE":                       "申請人上傳了過期的文件。",
	"FRONT_SIDE_MISSING":                    "文件的正面缺失。",
	"GRAPHIC_EDITOR":                        "文檔或其數據的亮度、對比度、內容等已更改。",
	"ID_INVALID":                            "識別個人身份的文件（如護照或身份證）無效。",
	"INCOMPATIBLE_LANGUAGE":                 "申請人應上傳其文件的譯文。",
	"INCOMPLETE_DOCUMENT":                   "文檔中缺少某些信息，或者部分信息可見。",
	"INCORRECT_SOCIAL_NUMBER":               "申請人提供的社會號碼（例如 SSN）不正確。",
	"PROBLEMATIC_APPLICANT_DATA":            "申請人數據與文件中的數據不符。",
	"REQUESTED_DATA_MISMATCH":               "提供的信息與從文檔中獲取的識別信息不匹配。",
	"SELFIE_WITH_PAPER":                     "申請人應上傳特殊的自拍照（例如，帶有文件和日期的自拍照）。",
	"LOW_QUALITY":                           "文件質量較低，無法得出明確的結論。",
	"NOT_ALL_CHECKS_COMPLETED":              "所有的檢查都還未完成。",
	"SCREENSHOTS":                           "申請人上傳了截圖。",
	"UNFILLED_ID":                           "申請人上傳的文件沒有簽名和蓋章。",
	"UNSATISFACTORY_PHOTOS":                 "照片存在問題，例如質量差或信息被掩蓋。",
	"UNSUITABLE_ENV":                        "在視頻識別通話中，申請人要麼不是獨自一人，要麼是看不見的。",
	"WRONG_ADDRESS":                         "文件中的地址與申請人輸入的地址不匹配。",
	"ADVERSE_MEDIA":                         "申請人在不利的媒體中被發現。",
	"AGE_REQUIREMENT_MISMATCH":              "不符合年齡要求（例如不能向 25 歲以下的人租車）。",
	"BLACKLIST":                             "申請人已被我們列入黑名單。",
	"BLOCKLIST":                             "申請人已被您列入黑名單。",
	"CHECK_UNAVAILABLE":                     "資料庫不可用。",
	"COMPROMISED_PERSONS":                   "申請人不符合妥協人政治。",
	"CRIMINAL":                              "申請人有違法行為。",
	"DB_DATA_MISMATCH":                      "資料不匹配；無法驗證配置檔案。",
	"DB_DATA_NOT_FOUND":                     "未找到任何資料；無法驗證配置檔案。",
	"DOCUMENT_TEMPLATE":                     "所提供的文件是從互聯網下載的模板。",
	"DUPLICATE":                             "該申請人已為此客戶創建，法規不允許重複。",
	"EXPERIENCE_REQUIREMENT_MISMATCH":       "經驗不足（例如駕駛經驗不夠）。",
	"FORGERY":                               "已進行偽造嘗試。",
	"FRAUDULENT_LIVENESS":                   "試圖繞過活性檢查。",
	"FRAUDULENT_PATTERNS":                   "檢測到欺詐行為。",
	"INCONSISTENT_PROFILE":                  "不同人的資料或文件被上傳給一名申請人。",
	"PEP":                                   "申請人屬於 PEP 類別。",
	"REGULATIONS_VIOLATIONS":                "違反規定。",
	"SANCTIONS":                             "該申請人被發現在制裁名單上。",
	"SELFIE_MISMATCH":                       "申請人照片（個人資料圖片）與所提供文件上的照片不符。",
	"SPAM":                                  "申請人被錯誤創建或只是垃圾郵件用戶（提供了不相關的圖像）。",
	"NOT_DOCUMENT":                          "提供的文件與驗證程序無關。",
	"THIRD_PARTY_INVOLVED":                  "申請人正在向第三方收費進行驗證。",
	"UNSUPPORTED_LANGUAGE":                  "視頻識別語言不受支持。",
	"WRONG_USER_REGION":                     "當某些地區/國家的申請人不允許註冊時。",
}
