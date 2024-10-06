part of '../ui_kit_part.dart';

class HeaderTextField extends StatelessWidget {
  final String? headerText;
  final TextStyle? headerStyle;
  final double? headerPadding;
  final CrossAxisAlignment? headerAlignment;
  final TextEditingController controller;
  final AutovalidateMode? autovalidateMode;
  final String? Function(String?)? validator;

  const HeaderTextField({
    super.key, 
    this.headerText,
    this.headerStyle,
    this.headerPadding,
    this.headerAlignment,
    required this.controller,
    this.autovalidateMode,
    this.validator,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: headerAlignment ?? CrossAxisAlignment.start,
      children: [
        if (headerText != null) ...[
          Text(
            headerText!,
            style: headerStyle,
            ),
          SizedBox(height: headerPadding ?? 16),
          ], 
        TextFormField(
          controller: controller,
          autovalidateMode: autovalidateMode,
          validator: validator,
        ),
      ],
    );
  }
}