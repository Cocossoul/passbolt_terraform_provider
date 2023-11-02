package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PasswordAndDescriptionResourceModel describes the resource data model.
type PasswordAndDescriptionResourceModel struct {
	Name           types.String `tfsdk:"name"`
	FolderParentId types.String `tfsdk:"folder_parent_id"`
	Username       types.String `tfsdk:"username"`
	URI            types.String `tfsdk:"uri"`
	Password       types.String `tfsdk:"password"`
	Description    types.String `tfsdk:"description"`
	Id             types.String `tfsdk:"id"`
}
