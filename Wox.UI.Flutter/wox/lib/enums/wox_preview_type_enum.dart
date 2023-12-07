typedef WoxPreviewType = String;

enum WoxPreviewTypeEnum {
  WOX_PREVIEW_TYPE_MARKDOWN("markdown", "markdown"),
  WOX_PREVIEW_TYPE_TEXT("text", "text"),
  WOX_PREVIEW_TYPE_IMAGE("image", "image"),
  WOX_PREVIEW_TYPE_URL("url", "url");

  final String code;
  final String value;

  const WoxPreviewTypeEnum(this.code, this.value);

  static String getValue(String code) => WoxPreviewTypeEnum.values.firstWhere((activity) => activity.code == code).value;
}