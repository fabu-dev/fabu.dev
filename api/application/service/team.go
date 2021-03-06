package service

import (
	"fabu.dev/api/application/controller/response"
	"fabu.dev/api/application/model"
	"fabu.dev/api/pkg/api"
	"fabu.dev/api/pkg/api/code"
	"fabu.dev/api/pkg/api/request"
	"fabu.dev/api/pkg/constant"
)

type Team struct {
}

func NewTeam() *Team {
	return &Team{}
}

// 获取会员的团队列表
func (s *Team) GetListByMember(memberId uint64) (*response.TeamList, *api.Error) {
	// 先获取会员所有的团队
	objTeamMember := model.NewTeamMember()
	teamIdSlice, err := objTeamMember.GetTeamId(memberId)
	if err != nil {
		return nil, err
	}

	// 获取团队信息
	objTeam := model.NewTeam()
	teamSlice, err := objTeam.GetListById(teamIdSlice)
	if err != nil {
		return nil, err
	}

	result := &response.TeamList{
		CountTeam:        len(teamSlice),
		CountApp:         0,
		CountAppDownload: 0,
		Team:             teamSlice,
	}

	return result, err
}

// 获取单个团队的成员信息
func (s *Team) GetMemberList(teamId uint64) ([]*model.TeamMemberInfo, *api.Error) {

	objTeamMember := model.NewTeamMember()
	teamMemberList, err := objTeamMember.GetListByTeamId(teamId)

	if len(teamMemberList) == 0 {
		return teamMemberList, nil
	}

	err = s.ApplyMember(teamMemberList)

	return teamMemberList, err
}

// 给团队成员列表应用用户名
func (s *Team) ApplyMember(teamMemberList []*model.TeamMemberInfo) *api.Error {
	memberIdList := make([]uint64, 0, len(teamMemberList))

	for _, teamMember := range teamMemberList {
		teamMember.RoleName = model.TeamRoleMap[teamMember.Role]
		memberIdList = append(memberIdList, teamMember.MemberId)
	}

	objMember := model.NewMember()
	memberList, err := objMember.GetListById(memberIdList)
	if err != nil {
		return err
	}

	for _, teamMember := range teamMemberList {
		for _, member := range memberList {
			if teamMember.MemberId == member.Id {
				teamMember.MemberName = member.Name
				teamMember.MemberAccount = member.Account
				teamMember.MemberEmail = member.Email
			}
		}
	}

	return nil
}

// 获取单个团队的成员信息
func (s *Team) DeleteMember(teamMemberId uint64) *api.Error {
	// TODO 删除之前需要做一些验证，eg：团队是否没人了，团队是否还有APP 是否是创建者
	objTeamMember := model.NewTeamMember()
	err := objTeamMember.Delete(teamMemberId)

	return err
}

// 邀请团队成员
func (s *Team) InviteMember(params *request.TeamMemberAddParams, operator *model.Operator) *api.Error {
	// 检查会员是否存在
	memberInfo, err := model.NewMember().GetDetailByEmail(params.Email)
	if err != nil {
		return err
	}

	if isIn, _ := model.NewTeamMember().IsInTeam(params.Id, memberInfo.Id); isIn {
		return api.NewError(code.ErrorTeamMemberExist, code.GetMessage(code.ErrorTeamMemberExist))
	}

	// 将会员加入团队
	teamMemberInfo := &model.TeamMemberInfo{
		TeamId:    params.Id,
		MemberId:  memberInfo.Id,
		Role:      params.Role,
		CreatedBy: operator.Account,
	}
	err = s.AddMember(teamMemberInfo)

	return err
}

// 创建团队
func (s *Team) Create(params *request.TeamCreateParams, operator *model.Operator) (*model.TeamInfo, *api.Error) {
	teamInfo := &model.TeamInfo{
		Owner:     operator.Id,
		Name:      params.Name,
		Status:    constant.StatusEnable,
		CreatedBy: operator.Account,
	}
	if err := s.CreateTeam(teamInfo); err != nil {
		return nil, err
	}

	teamMemberInfo := &model.TeamMemberInfo{
		TeamId:    teamInfo.Id,
		MemberId:  uint64(operator.Id),
		Role:      constant.TeamRoleCreator,
		CreatedBy: operator.Account,
	}
	err := s.AddMember(teamMemberInfo)

	return teamInfo, err
}

// 编辑团队
func (s *Team) Edit(params *request.TeamEditParams, operator *model.Operator) (*model.TeamInfo, *api.Error) {
	teamInfo := &model.TeamInfo{
		Id:        params.Id,
		Name:      params.Name,
		UpdatedBy: operator.Account,
	}

	objTeam := model.NewTeam()
	if err := objTeam.Edit(teamInfo); err != nil {
		return nil, err
	}

	return teamInfo, nil
}

// 删除团队
func (s *Team) Delete(params *request.TeamDeleteParams, operator *model.Operator) *api.Error {
	if isHas, _ := model.NewApp().HasByTeamId(params.Id); isHas {
		return api.NewError(code.ErrorTeamHasApp, code.GetMessage(code.ErrorTeamHasApp))
	}

	teamInfo := &model.TeamInfo{
		Id:        params.Id,
		Status:    constant.StatusDisable,
		UpdatedBy: operator.Account,
	}

	return s.DeleteTeam(teamInfo)
}

// 逻辑删除team表的记录
func (s *Team) DeleteTeam(teamInfo *model.TeamInfo) *api.Error {
	objTeam := model.NewTeam()

	if err := objTeam.Edit(teamInfo); err != nil {
		return err
	}

	return nil
}

// 保存数据到team表
func (s *Team) CreateTeam(teamInfo *model.TeamInfo) *api.Error {
	objTeam := model.NewTeam()

	return objTeam.Add(teamInfo)
}

// 保存数据到team_member表（将创建团队的人，加入到这个团队）
func (s *Team) AddMember(teamMemberInfo *model.TeamMemberInfo) *api.Error {
	objTeamMember := model.NewTeamMember()

	return objTeamMember.Add(teamMemberInfo)
}

// 获取团队详细信息
func (s *Team) GetInfoById(teamId uint64) (*model.TeamInfo, *api.Error) {
	objTeam := model.NewTeam()

	teamInfo, err := objTeam.GetInfoById(teamId)

	//s.ApplyPlatformName(appInfo)

	return teamInfo, err
}
