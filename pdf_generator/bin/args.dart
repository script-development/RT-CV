import 'dart:io';

import 'package:args/args.dart';
import 'package:pdf/pdf.dart';

import 'utils.dart';

class ArgsParser {
  ArgsParser(List<String> args) {
    ArgParser argsParser = ArgParser();
    argsParser.addFlag(
      'help',
      abbr: 'h',
      help: "Print this message",
    );
    argsParser.addFlag(
      'dummy',
      help:
          "Use dummy data, handy for working on this application. The dummy data is located in bin/cv.dart",
    );
    argsParser.addOption(
      'data',
      help:
          'input CV as json data (the structure of the CV should be the CV in /models/cv.go marshaled)',
    );
    argsParser.addOption(
      'header-color',
      help: 'set the backgorund color hex (#ffffff) of the main header',
      defaultsTo: '#4398a5',
    );
    argsParser.addOption(
      'sub-header-color',
      help: 'set the background color hex (#ffffff) of the sub headers',
      defaultsTo: '#ffe004',
    );
    argsParser.addOption(
      "logo-image-url",
      help: 'set the logo image url, leave empty for no logo',
    );
    argsParser.addOption(
      "company-name",
      help: 'set the company name',
    );
    argsParser.addOption(
      "company-address",
      help: 'set the company address section',
    );
    argsParser.addOption(
      "font-regular",
      help:
          'set the font to use, for a list of fonts look in bin/fonts.dart > _fontFilesMap',
      defaultsTo: 'OpenSans',
    );
    argsParser.addOption(
      "font-bold",
      help: 'set the font to use',
      defaultsTo: 'OpenSans',
    );
    argsParser.addOption(
      'style',
      help:
          'set the style of the document based an a set list of styles, for a list of styles look in bin/args.dart > LayoutStyle',
      defaultsTo: 'style_1',
    );
    argsParser.addOption(
      'out',
      abbr: 'o',
      defaultsTo: 'example.pdf',
      help: "to where should we write the output file",
    );

    try {
      this.argResult = argsParser.parse(args);
    } catch (e) {
      print('unable to parse args, error: ${e}');
      exit(1);
    }
    if (argResult['help'] == true) {
      print(argsParser.usage);
      exit(0);
    }
  }

  ArgsParser.fromArgResults(this.argResult);

  late final ArgResults argResult;

  /// data is the cv data (should be encoded as json)
  /// The layout of this data is based of the main project's CV struct in /models/cv.go
  String get data => argResult['data'] ?? '';

  /// should dummy data be used to generate the document
  bool get dummy => argResult['dummy'];

  /// logoImageUrl is a logo place on the bottom of the document
  /// This logo is fetched from the internet
  String? get logoImageUrl => argResult['logo-image-url'];

  /// fontRegular is the regular font to use
  String get fontRegular => argResult['font-regular'];

  /// fontBold is font used as bold font thus used for the headers
  String get fontBold => argResult['font-bold'];

  /// The company name placed at the bottom of the file
  String? get companyName => argResult['company-name'];

  /// The company address placed at the bottom of the file
  String? get companyAddress => argResult['company-address'];

  /// The style of the document
  Style get style => Style(
        layoutStyle: pdfStyleFromString(argResult['style']),
        headerBackgroundColor: PdfColor.fromHex(argResult['header-color']),
        subHeaderBackgroundColor:
            PdfColor.fromHex(argResult['sub-header-color']),
      );

  /// The output file name
  String get out => argResult['out'];
}

class Style {
  Style({
    required this.layoutStyle,
    required this.headerBackgroundColor,
    required this.subHeaderBackgroundColor,
  }) {
    this.headerTextColor = getTextColorFromBg(headerBackgroundColor);
    this.subHeaderTextColor = getTextColorFromBg(subHeaderBackgroundColor);
  }

  final LayoutStyle layoutStyle;

  final PdfColor headerBackgroundColor;
  late final PdfColor headerTextColor;

  final PdfColor subHeaderBackgroundColor;
  late final PdfColor subHeaderTextColor;
}

enum LayoutStyle {
  style_1,
  style_2,
  style_3,
}

LayoutStyle pdfStyleFromString(String style) {
  switch (style) {
    case 'style_1':
      return LayoutStyle.style_1;
    case 'style_2':
      return LayoutStyle.style_2;
    case 'style_3':
      return LayoutStyle.style_3;
    default:
      print('unknown style: ${style}');
      return LayoutStyle.style_1;
  }
}
