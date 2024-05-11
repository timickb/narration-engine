package mocks

//go:generate mockgen -package=mocks -source=../domain/interfaces.go -destination=TransactorMock.go -mock_names=Transactor=TransactorMock
//go:generate mockgen -package=mocks -source=../core/interfaces.go -destination=CoreMocks.go -mock_names=AsyncWorkerConfig=AsyncWorkerConfigMock,HandlerAdapter=HandlerAdapterMock,InstanceRepository=CoreInstanceRepositoryMock,TransitionRepository=CoreTransitionRepositoryMock
//go:generate mockgen -package=mocks -source=../usecase/interfaces.go -destination=UsecaseMocks.go -mock_names=InstanceRepository=UsecaseInstanceRepositoryMock,PendingEventRepository=UsecasePendingEventRepositoryMock,Config=UsecaseConfigMock
