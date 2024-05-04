package parser

import "github.com/timickb/narration-engine/internal/domain"

var (
	// Состояния для тестового сценария 1
	testState11 = &domain.State{
		Name:    "Состояние1",
		Handler: "service.worker1",
		Params:  nil,
	}
	testState12 = &domain.State{
		Name:    "Состояние2",
		Handler: "service.worker2",
		Params:  nil,
	}
	testState13 = &domain.State{
		Name:    "Состояние3",
		Handler: "",
		Params:  nil,
	}

	// Состояния для тестового сценария 2
	testState2Fetch = &domain.State{
		Name:    "FetchSomething",
		Handler: "service1.fetch_something",
	}
	testState2Good = &domain.State{
		Name:    "SomethingIsGood",
		Handler: "service2.do_if_good",
	}
	testState2Send = &domain.State{
		Name:    "SendNotification",
		Handler: "service2.send_notification",
	}
	testState2Bad = &domain.State{
		Name:    "SomethingIsBad",
		Handler: "service2.do_if_bad",
	}

	// Состояния для тестового сценария 3
	testState3FetchOrderInfo = &domain.State{
		Name:    "FetchOrderInfo",
		Handler: "order_service.fetch_order",
		Params: map[string]domain.StateParamValue{
			"order_id": {
				Value:       "123",
				FromContext: false,
			},
			"user_id": {
				Value:       "456",
				FromContext: false,
			},
		},
	}
	testState3CreateInvoice = &domain.State{
		Name:    "CreateInvoice",
		Handler: "payment_service.create_invoice",
	}
	testState3SendSuccessEmail = &domain.State{
		Name:    "SendSuccessEmail",
		Handler: "notification_service.send_email",
		Params: map[string]domain.StateParamValue{
			"data": {
				Value:       "ctx.success_email",
				FromContext: true,
			},
		},
	}
	testState3SendFailEmail = &domain.State{
		Name:    "SendFailEmail",
		Handler: "notification_service.send_email",
		Params: map[string]domain.StateParamValue{
			"data": {
				Value:       "ctx.fail_email",
				FromContext: true,
			},
		},
	}

	// Состояния для тестового сценария модерации публикаций

)
