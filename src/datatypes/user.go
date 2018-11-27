package datatypes

type MsgrUserJSON struct {
	MsgrUserID int     `json:"msgrUserId"`
	Username   *string `json:"username"`
	Password   *string `json:"password"`
}
