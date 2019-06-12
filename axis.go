package diamond_client

type Axis struct {
    group  string
    dataId string
}

func NewAxis(group string, dataId string) *Axis {
    if "" == dataId {
        panic("blank dataId")
    }
    if "" == group {
        group = DefaultGroupName
    }
    return &Axis{group: group, dataId: dataId}
}

func NewAxisOfDefaultGroup(dataId string) *Axis {
    return NewAxis("", dataId)
}

func (axis *Axis) GetGroup() string {
    return axis.group
}

func (axis *Axis) GetDataId() string {
    return axis.dataId
}

func (axis *Axis) String() string {
    return "DiamondAxis{dataId=" + axis.dataId + ", group=" + axis.group + "}"
}
