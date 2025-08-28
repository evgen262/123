package employees

import (
	"time"

	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	managementv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/management/v1"
	positionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/position/v1"
	subdivisionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/subdivision/v1"
	workplacev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/workplace/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

type mapperEmployees struct {
	tu timeUtils.TimeUtils
}

func NewMapperEmployees(tu timeUtils.TimeUtils) *mapperEmployees {
	return &mapperEmployees{tu: tu}
}

func (m *mapperEmployees) EmployeeToEntity(employee *employeev1.CompositeEmployee) *entityEmployee.Employee {
	if employee == nil {
		return nil
	}

	var (
		orgGlobalID *string
		orgFullName *string
	)

	if employee.GetOrganization().GetGlobalId().GetValue() != "" {
		_g := employee.GetOrganization().GetGlobalId().GetValue()
		orgGlobalID = &_g
	}

	if employee.GetOrganization().GetFullName() != "" {
		_n := employee.GetOrganization().GetFullName()
		orgFullName = &_n
	}

	id, err := uuid.Parse(employee.GetId())
	if err != nil {
		id = uuid.Nil
	}

	orgID, err := uuid.Parse(employee.GetOrganization().GetId())
	if err != nil {
		orgID = uuid.Nil
	}

	var (
		mainProduct *entityEmployee.Product
		products    []*entityEmployee.Product
	)
	for _, product := range employee.GetProducts() {
		p := m.ProductToEntity(product)
		if p.IsMain {
			mainProduct = p
		}
		products = append(products, p)
	}
	if mainProduct == nil && len(products) > 0 {
		mainProduct = products[0]
	}

	e := &entityEmployee.Employee{
		ID: id,
		Portal: entityEmployee.Portal{
			ID:       int(employee.GetPortal().GetId()),
			Name:     employee.GetPortal().GetName(),
			IsActive: employee.GetPortal().GetIsActive(),
		},
		Person:         m.EmployeePersonToEntity(employee.GetPerson()),
		Number:         employee.GetNumber(),
		FullName:       employee.GetFullName(),
		EmploymentType: employee.GetEmploymentType(),
		Phones:         m.PhonesToEntity(employee.GetPhones()),
		Email:          employee.GetEmail(),
		Rate:           employee.GetRate(),
		Organization: entityEmployee.Organization{
			ID:              orgID,
			GlobalID:        orgGlobalID,
			Name:            employee.GetOrganization().GetName(),
			ShortName:       employee.GetOrganization().GetShortName(),
			FullName:        orgFullName,
			INN:             employee.GetOrganization().GetInn(),
			IsActive:        employee.GetOrganization().GetIsActive(),
			IsLiquidated:    employee.GetOrganization().GetIsLiquidated(),
			CreatedTime:     m.tu.TimestampToTime(employee.GetOrganization().GetCreatedTime()),
			LiquidationTime: m.tu.TimestampToTime(employee.GetOrganization().GetLiquidationTime()),
		},
		MainProduct:     mainProduct,
		Products:        products,
		StaffPosition:   m.StaffPositionToEntity(employee),
		Managements:     m.ManagmentsToEntity(employee.GetManagements()),
		SubdivisionTree: m.SubdivisionTreeToEntity(employee.GetSubdivisionTree()),
		ManagementTree:  m.ManagementTreeToEntity(employee.GetManagementTree()),
		Workplace:       m.WorkplaceToEntity(employee.GetWorkplace()),
		History:         m.HistoriesToEntity(employee.GetHistory()),
		IsFired:         employee.GetIsFired(),
		CreateTime:      m.tu.TimestampToTime(employee.GetCreateTime()),
		UpdateTime:      m.tu.TimestampToTime(employee.GetUpdateTime()),
	}

	empPhotoID, err := uuid.Parse(employee.GetPhoto())
	if err == nil {
		e.ImageID = empPhotoID
	}

	if d := employee.GetEmploymentTime(); d != nil {
		t := d.AsTime()
		e.DateOfEmployment = (*timeUtils.Date)(&t)
	}

	return e
}

func (m *mapperEmployees) EmployeePersonToEntity(personPb *employeev1.CompositeEmployee_Person) entityEmployee.Person {
	if personPb == nil {
		return entityEmployee.Person{}
	}

	id, err := uuid.Parse(personPb.GetId())
	if err != nil {
		id = uuid.Nil
	}

	person := entityEmployee.Person{
		ID:         id,
		ExtID:      personPb.GetExtId(),
		CloudID:    personPb.GetCloudId(),
		LastName:   personPb.GetLastName(),
		FirstName:  personPb.GetFirstName(),
		MiddleName: personPb.GetMiddleName(),
		INN:        personPb.GetInn(),
		SNILS:      personPb.GetSnils(),
		Gender:     m.GenderToEntity(personPb.GetGender()),
		Phone:      personPb.GetPhone(),
		Socials:    m.CompositePersonSocialsToEntity(personPb.GetSocials()),
		IsActive:   personPb.GetIsActive(),
	}

	if birthday := personPb.GetBirthday(); birthday != "" {
		if date, paresErr := timeUtils.ParseDate(birthday); paresErr == nil {
			person.Birthday = date
		}
	}

	return person
}

func (m *mapperEmployees) GenderToEntity(gender employeev1.GenderType) entity.Gender {
	switch gender {
	case employeev1.GenderType_GENDER_TYPE_MAN:
		return entity.GenderMale
	case employeev1.GenderType_GENDER_TYPE_WOMAN:
		return entity.GenderFemale
	default:
		return entity.GenderInvalid
	}
}

func (m *mapperEmployees) CompositePersonSocialsToEntity(socials *employeev1.CompositeEmployee_Person_Socials) *entityEmployee.Socials {
	if socials == nil {
		return nil
	}

	var chatID *string

	if socials.GetChatId() != "" {
		cID := socials.GetChatId()
		chatID = &cID
	}

	return &entityEmployee.Socials{
		Telegram: socials.GetTelegram(),
		ChatID:   chatID,
	}
}

func (m *mapperEmployees) StaffPositionToEntity(employee *employeev1.CompositeEmployee) *entityEmployee.StaffPosition {
	if employee.GetStaffPosition() == nil {
		return nil
	}

	id, err := uuid.Parse(employee.GetStaffPosition().GetId())
	if err != nil {
		id = uuid.Nil
	}

	orgID, err := uuid.Parse(employee.GetStaffPosition().GetOrganizationId())
	if err != nil {
		orgID = uuid.Nil
	}

	var createD *timeUtils.Date
	_createD, _ := time.Parse(time.DateOnly, employee.GetStaffPosition().GetCreateDate())
	if !_createD.IsZero() {
		createD = (*timeUtils.Date)(&_createD)
	}

	var closeD *timeUtils.Date
	_closeD, _ := time.Parse(time.DateOnly, employee.GetStaffPosition().GetCloseDate())
	if !_closeD.IsZero() {
		closeD = (*timeUtils.Date)(&_closeD)
	}

	return &entityEmployee.StaffPosition{
		ID: id,
		Portal: entityEmployee.Portal{
			ID: int(employee.StaffPosition.GetPortalId()),
		},
		Organization: entityEmployee.Organization{
			ID: orgID,
		},
		Subdivision:   m.SubdivisionToEntity(employee.GetSubdivision()),
		Position:      m.PositionToEntity(employee.GetPosition()),
		ResponsibleID: employee.GetStaffPosition().GetPositionId(),
		Name:          employee.GetStaffPosition().GetName(),
		RateNumbers:   int(employee.GetStaffPosition().GetRateNumbers()),
		CreateDate:    createD,
		CloseDate:     closeD,
		CreatedTime:   m.tu.TimestampToTime(employee.GetStaffPosition().GetCreatedTime()),
		UpdatedTime:   m.tu.TimestampToTime(employee.GetStaffPosition().GetUpdatedTime()),
	}
}

func (m *mapperEmployees) PositionToEntity(position *positionv1.Position) *entityEmployee.Position {
	if position == nil {
		return nil
	}

	id, err := uuid.Parse(position.GetId())
	if err != nil {
		id = uuid.Nil
	}

	return &entityEmployee.Position{
		ID:   id,
		Name: position.GetName(),
	}
}

func (m *mapperEmployees) SubdivisionToEntity(subdivision *subdivisionv1.Subdivision) *entityEmployee.Subdivision {
	if subdivision == nil {
		return nil
	}

	id, err := uuid.Parse(subdivision.GetId())
	if err != nil {
		id = uuid.Nil
	}

	return &entityEmployee.Subdivision{
		ID:        id,
		Name:      subdivision.GetName(),
		ParentID:  subdivision.GetParentId(),
		Sort:      int(subdivision.GetSort()),
		IsDeleted: subdivision.GetIsDeleted(),
	}
}

func (m mapperEmployees) PhonesToEntity(phones *employeev1.Phones) *entityEmployee.Phones {
	if phones == nil {
		return nil
	}

	var (
		wp *string
		ep *string
	)

	if phones.GetWorkNumber() != "" {
		_p := phones.GetWorkNumber()
		wp = &_p
	}

	if phones.GetExtensionNumber() != "" {
		_p := phones.GetExtensionNumber()
		ep = &_p
	}

	return &entityEmployee.Phones{
		WorkNumber:      wp,
		ExtensionNumber: ep,
	}
}

func (m *mapperEmployees) ManagementToEntity(management *managementv1.Management) *entityEmployee.Management {
	if management == nil {
		return nil
	}

	id, err := uuid.Parse(management.GetId())
	if err != nil {
		id = uuid.Nil
	}

	_m := &entityEmployee.Management{
		ID: id,
		Portal: entityEmployee.Portal{
			ID: int(management.GetPortalId()),
		},
		RoleID:      management.GetRoleId(),
		RoleName:    management.GetRoleName(),
		RoleType:    management.GetRoleType(),
		EmployeeID:  management.GetEmployeeId(),
		ProductID:   management.GetProductId(),
		IsMain:      management.GetIsMain(),
		IsDeleted:   management.GetIsDeleted(),
		CreatedTime: m.tu.TimestampToTime(management.GetCreatedTime()),
		UpdatedTime: m.tu.TimestampToTime(management.GetUpdatedTime()),
	}

	if management.GetParentId() != uuid.Nil.String() {
		_m.ParentID = management.GetParentId()
	}

	return _m
}

func (m mapperEmployees) ManagmentsToEntity(managments []*managementv1.Management) []*entityEmployee.Management {
	if managments == nil {
		return nil
	}

	arr := make([]*entityEmployee.Management, 0, len(managments))
	for _, management := range managments {
		_management := m.ManagementToEntity(management)
		if _management != nil {
			arr = append(arr, _management)
		}
	}
	return arr
}

func (m *mapperEmployees) WorkplaceToEntity(workplace *workplacev1.Workplace) *entityEmployee.Workplace {
	if workplace == nil {
		return nil
	}
	return &entityEmployee.Workplace{
		Address:       workplace.GetBuilding().GetAddress(),
		Floor:         int(workplace.GetFloor().GetNumber()),
		Number:        workplace.GetWorkplaceNumber().GetNumber(),
		CabinetNumber: workplace.GetCabinet().GetNumber(),
	}
}

func (m *mapperEmployees) HistoryToEntity(history *employeev1.History) *entityEmployee.History {
	if history == nil {
		return nil
	}

	id, err := uuid.Parse(history.GetId())
	if err != nil {
		id = uuid.Nil
	}

	h := &entityEmployee.History{
		ID:             id,
		EventType:      m.HistoryEventTypeToEntity(history.EventType),
		EmployeeNumber: history.EmployeeNumber,
		EmploymentType: history.EmploymentType,
		Rate:           history.Rate,
	}

	t := history.GetEventTime().AsTime()
	if !t.IsZero() {
		h.EventTime = &t
	}

	return h
}

func (m *mapperEmployees) HistoriesToEntity(histories []*employeev1.History) []*entityEmployee.History {
	if histories == nil {
		return nil
	}

	arr := make([]*entityEmployee.History, 0, len(histories))
	for _, history := range histories {
		_history := m.HistoryToEntity(history)
		if _history != nil {
			arr = append(arr, _history)
		}
	}
	return arr
}

func (m *mapperEmployees) SubdivisionTreeToEntity(tree *subdivisionv1.SubdivisionTree) *entityEmployee.SubdivisionTree {
	if tree == nil {
		return nil
	}

	children := make([]*entityEmployee.SubdivisionTree, 0, len(tree.Children))
	for _, child := range tree.Children {
		convertedChild := m.SubdivisionTreeToEntity(child) // Рекурсивный вызов
		if convertedChild != nil {
			children = append(children, convertedChild)
		}
	}

	id, err := uuid.Parse(tree.Id)
	if err != nil {
		id = uuid.Nil
	}

	parentID, err := uuid.Parse(tree.ParentId)
	if err != nil {
		parentID = uuid.Nil
	}

	return &entityEmployee.SubdivisionTree{
		ID:                   id,
		Name:                 tree.Name,
		ManagementPositionID: tree.ManagementPositionId,
		ParentID:             parentID.String(),
		Children:             children,
		Sort:                 int(tree.Sort),
		IsDeleted:            tree.IsDeleted,
	}
}

func (m *mapperEmployees) SubdivisionTreesToEntity(trees []*subdivisionv1.SubdivisionTree) []*entityEmployee.SubdivisionTree {
	if trees == nil {
		return nil
	}

	arr := make([]*entityEmployee.SubdivisionTree, 0, len(trees))
	for _, tree := range trees {
		convertedTree := m.SubdivisionTreeToEntity(tree)
		if convertedTree != nil {
			arr = append(arr, convertedTree)
		}
	}
	return arr
}

func (m *mapperEmployees) ManagementTreeToEntity(tree *managementv1.ManagementTree) *entityEmployee.ManagementTree {
	if tree == nil {
		return nil
	}

	children := make([]*entityEmployee.ManagementTree, 0, len(tree.Children))
	for _, child := range tree.Children {
		convertedChild := m.ManagementTreeToEntity(child) // Рекурсивный вызов
		if convertedChild != nil {
			children = append(children, convertedChild)
		}
	}

	id, err := uuid.Parse(tree.Id)
	if err != nil {
		id = uuid.Nil
	}

	parentID, err := uuid.Parse(tree.ParentId)
	if err != nil {
		parentID = uuid.Nil
	}

	return &entityEmployee.ManagementTree{
		ID:        id,
		Name:      tree.Name,
		RoleName:  tree.RoleName,
		RoleType:  tree.RoleType,
		ParentID:  parentID.String(),
		Children:  children,
		Sort:      int(tree.Sort),
		IsDeleted: tree.IsDeleted,
	}
}

func (m *mapperEmployees) ManagementTreesToEntity(trees []*managementv1.ManagementTree) []*entityEmployee.ManagementTree {
	if trees == nil {
		return nil
	}

	arr := make([]*entityEmployee.ManagementTree, 0, len(trees))
	for _, tree := range trees {
		convertedTree := m.ManagementTreeToEntity(tree)
		if convertedTree != nil {
			arr = append(arr, convertedTree)
		}
	}
	return arr
}

func (m *mapperEmployees) ProductToEntity(productPb *employeev1.Product) *entityEmployee.Product {
	p := &entityEmployee.Product{
		ID:            productPb.GetId(),
		PortalID:      productPb.GetPortalId(),
		ClusterID:     productPb.GetClusterId(),
		ClusterName:   productPb.GetClusterName(),
		FullName:      productPb.GetFullName(),
		ShortName:     productPb.GetShortName(),
		Type:          productPb.GetType(),
		ResponsibleID: productPb.GetResponsibleId(),
		Tutor:         productPb.GetTutorId(),
		IsMain:        productPb.GetIsMain(),
	}

	iconID, err := uuid.Parse(productPb.GetIcon())
	if err == nil {
		p.IconID = iconID
	}

	return p
}

func (m mapperEmployees) HistoryEventTypeToEntity(e employeev1.HistoryEventType) entityEmployee.OperationType {
	switch e {
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_HIRING:
		return entityEmployee.OperationTypeHiring
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_TRANSFER:
		return entityEmployee.OperationTypeTransfer
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_DISMISSAL:
		return entityEmployee.OperationTypeDismissal
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_WORKING_CONDITIONS:
		return entityEmployee.OperationTypeWorkingConditions
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_NAME_CHANGE:
		return entityEmployee.OperationTypeNameChange
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_CHILDCARE_LEAVE:
		return entityEmployee.OperationTypeChildcareLeave
	case employeev1.HistoryEventType_HISTORY_EVENT_TYPE_MATERNITY_LEAVE:
		return entityEmployee.OperationTypeMaternityLeave
	default:
		return entityEmployee.OperationTypeInvalid
	}
}
