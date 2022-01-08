import 'dart:io';
import 'dart:typed_data';

import 'package:pdf/widgets.dart';

String? formatDate(DateTime? input) {
  if (input == null) return null;

  String year = input.year.toString();
  String month = input.month.toString().padLeft(2, '0');
  String day = input.day.toString().padLeft(2, '0');

  return "${year}-${month}-${day}";
}

String? formatDateTime(DateTime? input) {
  if (input == null) return null;

  String hour = input.hour.toString().padLeft(2, '0');
  String minute = input.minute.toString().padLeft(2, '0');
  String second = input.second.toString().padLeft(2, '0');

  return "${formatDate(input)!} ${hour}:${minute}:${second}";
}

Future<Font> loadFont(String file) async {
  final File fontFile = File(file);
  Uint8List data = await fontFile.readAsBytesSync();
  return Font.ttf(ByteData.view(data.buffer));
}
