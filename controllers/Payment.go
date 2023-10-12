package controllers

import (
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
)

type ControllerPayment struct {
}

func (cp *ControllerPayment) GetDebit(context *types.RequestContext) {
	context.Response.Data = append(context.Response.Data, webSocketModels.Payment{
		Period:      "01.2022",
		Date:        "31.01.2022",
		Allocate:    2332327,
		Payer:       "ПАО СБЕРБАНК//МАМЧУЕВ АЛЬБЕРТ ШАУХАЛОВИЧ//1316557498137//115211,Рататата",
		Destination: "СЧЕТ № 22-1219 ОТ 07.09.2022, НДС не облагается.",
		Scans:       webSocketModels.PaymentScans{},
	})
	context.Response.Data = append(context.Response.Data, webSocketModels.Payment{
		Period: "01.2022",
		Date:   "28.01.2022",
		Split: append([]webSocketModels.PaymentSplit{}, webSocketModels.PaymentSplit{
			Order:         32456,
			Employee:      1,
			Object:        1,
			TypeOfPayment: 1,
			Credit:        7,
		}),
		Scans: webSocketModels.PaymentScans{
			MediaTypes: 1,
			Documents:  1,
		},
		Allocate:    2332327,
		Payer:       "ПАО СБЕРБАНК//МАМЧУЕВ АЛЬБЕРТ ШАУХАЛОВИЧ//1316557498137//115211,Рататата",
		Destination: "СЧЕТ № 22-1219 ОТ 07.09.2022, НДС не облагается.",
	})
	context.Response.Data = append(context.Response.Data, webSocketModels.Payment{
		Period: "01.2022",
		Date:   "21.01.2022",
		Split: append(append([]webSocketModels.PaymentSplit{}, webSocketModels.PaymentSplit{
			Order:         32456,
			Employee:      0,
			Object:        1,
			TypeOfPayment: 1,
			Credit:        1,
		}), webSocketModels.PaymentSplit{
			Order:         32456,
			Employee:      0,
			Object:        1,
			TypeOfPayment: 1,
			Credit:        1,
		}),
		Scans: webSocketModels.PaymentScans{
			Documents: 1,
		},
		Allocate: 110000000,
		Payer: "Иванов Вячеслав Александрович (ИП) р/с 40802810501830000182 в Ф-Л СИБИРСКИЙ ПАО БАНК " +
			"\"ФК ОТКРЫТИЕ\" г Новосибирск",
		Destination: "ОПЛ. ЗА УСЛУГИ ПО ОРГАНИЗАЦИИ ЭФИРА ЛОКАЛЬНОЙ СЕТИ РАДИОВЕЩАНИЯ ЗА ИЮЛЬ 2022 Г. ПО ДОГОВОРУ" +
			" № А1-74 ОТ 01.03.2022 Г.",
	})
	context.Response.Data = append(context.Response.Data, webSocketModels.Payment{
		Period: "01.2022",
		Date:   "08.01.2022",
		Split: append([]webSocketModels.PaymentSplit{}, webSocketModels.PaymentSplit{
			Order:         21452,
			Employee:      137,
			Object:        64,
			TypeOfPayment: 1,
			Credit:        12791.80,
		}),
		Scans: webSocketModels.PaymentScans{
			Documents:  1,
			MediaTypes: 1,
			Payments:   1,
		},
		Allocate:    12791.80,
		Payer:       "ООО \"Эс Би Си Медиа\"",
		Destination: "ОПЛ. ПО СЧЕТУ №22-1121 ОТ 29.08.2022Г. ЗА АУДИОРОЛИК. НДС НЕ ОБЛАГ.",
	})
}
