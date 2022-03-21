import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'cv.dart';
import 'layout.dart';
import 'utils.dart';

class ClientInfo extends StatelessWidget {
  ClientInfo({
    required this.personalInfo,
    this.driversLicenses,
  });

  List<Widget> children = [];
  final PersonalDetails personalInfo;
  final List<String>? driversLicenses;

  final TextStyle labelStyle = TextStyle(
    fontSize: 10,
    color: PdfColors.grey800,
  );
  final TextStyle valueStyle = TextStyle(
    fontSize: 10,
    color: PdfColors.black,
  );

  tryAddToList(String label, String? value) {
    if (value != null) {
      children.add(
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(label + ": ", style: labelStyle),
            Text(
              value,
              overflow: TextOverflow.clip,
              style: valueStyle,
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(Context context) {
    children = [];
    tryAddToList("Geboortedatum", formatDate(personalInfo.dob));
    tryAddToList("E-mail", personalInfo.email);
    tryAddToList("Telefoon", personalInfo.phoneNumber);
    if (driversLicenses != null) {
      switch (driversLicenses!.length) {
        case 0:
          // Do not add the drivers licenses
          break;
        case 1:
          tryAddToList("Rijbewijs", driversLicenses![0]);
          break;
        default:
          tryAddToList("Rijbewijzen", driversLicenses!.join(", "));
      }
    }

    final EdgeInsets padding = layoutBlockBasePadding.copyWith(
      top: PdfPageFormat.cm * 1.5,
      bottom: PdfPageFormat.cm / 2,
    );

    if (!personalInfo.hasAddress) {
      if (personalInfo.zip != null) {
        String? postalCodePlace = guessPostalCodePlace(personalInfo.zip!);

        if (postalCodePlace != null)
          tryAddToList(
              "Postcode", "${personalInfo.zip} (regio ${postalCodePlace})");
        else
          tryAddToList("Postcode", personalInfo.zip);
      }

      return Padding(
        padding: padding,
        child: Wrap(children: children, spacing: 10),
      );
    }

    return Padding(
      padding: padding,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            constraints: BoxConstraints(
              minWidth: 150,
            ),
            child: Padding(
              padding: EdgeInsets.only(right: 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text("Plaats: ", style: labelStyle),
                      Text(personalInfo.city!, style: valueStyle),
                    ],
                  ),
                  Row(
                    children: [
                      Text("Adres: ", style: labelStyle),
                      Text(
                          "${personalInfo.streetName} ${personalInfo.houseNumber} ${personalInfo.houseNumberSuffix ?? ""}",
                          style: valueStyle),
                    ],
                  ),
                  Row(
                    children: [
                      Text("Postcode: ", style: labelStyle),
                      Text(personalInfo.zip!, style: valueStyle),
                    ],
                  ),
                ],
              ),
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: children,
          ),
        ],
      ),
    );
  }
}

class WorkExpWidget extends StatelessWidget {
  WorkExpWidget(WorkExperience this.exp);

  final WorkExperience exp;

  @override
  Widget build(Context context) {
    return ListEntry(exp.profession ?? '??',
        company: exp.employer,
        description: exp.description,
        from: exp.startDate,
        to: exp.endDate,
        padding: EdgeInsets.only(top: 15));
  }
}

class EducationWidget extends StatelessWidget {
  EducationWidget(Education this.education);

  final Education education;

  @override
  Widget build(Context context) {
    return ListEntry(
      education.name,
      company: education.institute,
      from: education.startDate,
      to: education.endDate,
      description: education.description,
      padding: EdgeInsets.only(top: 15),
    );
  }
}

class ListEntry extends StatelessWidget {
  ListEntry(
    this.title, {
    this.description,
    this.company,
    this.from,
    this.to,
    EdgeInsets? padding,
  }) : this.padding = padding ?? EdgeInsets.only(top: 5);

  final String title;
  final String? company;
  final String? description;
  final DateTime? from;
  final DateTime? to;
  final EdgeInsets padding;

  @override
  Widget build(Context context) {
    List<Widget> children = [
      Flexible(
        child: SafeText(
          title,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontSize: 10,
            fontWeight: FontWeight.bold,
            lineSpacing: 2,
          ),
        ),
      ),
    ];

    TextStyle contentStyle = TextStyle(
      fontSize: 10,
      color: PdfColors.grey800,
      lineSpacing: 2,
    );
    TextStyle dateStyle = TextStyle(
      fontSize: 8,
      color: PdfColors.grey600,
      lineSpacing: 2,
    );

    if (company != null && company!.isNotEmpty) {
      children.add(
        Row(children: [
          SafeText(
            company!,
            overflow: TextOverflow.clip,
            style: contentStyle,
          ),
        ]),
      );
    }

    String? fromStr = formatDate(from);
    String? toStr = formatDate(to);
    if (fromStr != null || toStr != null) {
      if (toStr == null || fromStr == null) {
        children.add(
          Row(children: [
            Text(fromStr ?? toStr ?? '??', style: dateStyle),
          ]),
        );
      } else {
        children.add(
          Row(children: [
            Text(fromStr, style: dateStyle),
            Text(" - ", style: dateStyle),
            Text(toStr, style: dateStyle),
          ]),
        );
      }
    }

    if (description != null && description!.isNotEmpty) {
      children.add(Flexible(
        child: SafeText(
          (description!.length > 300)
              ? description!.substring(0, 300) + '..'
              : description!,
          overflow: TextOverflow.clip,
          style: contentStyle,
        ),
      ));
    }

    return Padding(
      padding: const EdgeInsets.only(top: 5),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: children,
      ),
    );
  }
}

// SafeText makes it less likely to see boxes with "x" inside of them in the result pdf
// This happends most of the time because characters are not supported by the font
// This widget filters the text before passing it to the Text widget
class SafeText extends Text {
  SafeText(
    String text, {
    TextStyle? style,
    TextAlign? textAlign,
    TextDirection? textDirection,
    bool? softWrap,
    bool tightBounds = false,
    double textScaleFactor = 1.0,
    int? maxLines,
    TextOverflow? overflow,
  }) : super(
          _makeTextSafe(text),
          style: style,
          textAlign: textAlign,
          textDirection: textDirection,
          softWrap: softWrap,
          tightBounds: tightBounds,
          textScaleFactor: textScaleFactor,
          maxLines: maxLines,
          overflow: overflow,
        );
}

var _a = "a".codeUnitAt(0);
var _z = "z".codeUnitAt(0);
var _A = "A".codeUnitAt(0);
var _Z = "Z".codeUnitAt(0);
var _0 = "0".codeUnitAt(0);
var _9 = "9".codeUnitAt(0);
var _dash = "-".codeUnitAt(0);

var spaceChars = Runes(" \n");
var otherSpaceChars = Runes("\t\r\f\v");

String _makeTextSafe(String text) {
  List<int> safeChars = [];
  for (int c in text.runes) {
    if (c > 0x80) {
      // There are a lot of character points checked here,
      // These can be found on: https://www.compart.com/en/unicode/block

      // Latin-1 Supplement
      if (c > 0x00A0 && c < 0x00FF) {
        safeChars.add(c);
        continue;
      }

      // Latin Extended-A
      if (c >= 0x0100 && c <= 0x017F) {
        safeChars.add(c);
        continue;
      }

      // Latin Extended-A
      if (c >= 0x0180 && c <= 0x024F) {
        safeChars.add(c);
        continue;
      }

      // IPA Extensions
      if (c >= 0x0250 && c <= 0x02AF) {
        safeChars.add(c);
        continue;
      }

      // Space modifier letters
      if (c >= 0x02B0 && c <= 0x02FF) {
        safeChars.add(c);
        continue;
      }

      // Latin Extended Additional
      if (c >= 0x1E00 && c <= 0x1EFF) {
        safeChars.add(c);
        continue;
      }

      // General Punctuation
      if ((c >= 0x2000 && c <= 0x200F) || (c >= 0x2028 && c <= 0x202F)) {
        // These are space like characters, add a space
        if (safeChars.isNotEmpty && !spaceChars.contains(safeChars.last)) {
          safeChars.add(spaceChars.first);
        }
        continue;
      }
      if (c >= 0x2010 && c <= 0x2015) {
        // These are dash like characters, add a normal dash
        safeChars.add(_dash);
        continue;
      }
      if (c >= 0x2016 && c < 0x205E) {
        // Remainder of the General Punctuation block
        safeChars.add(c);
        continue;
      }

      // Currency Symbols
      if (c > 0x20A0 && c < 0x20CF) {
        safeChars.add(c);
        continue;
      }

      // Miscellaneous Technical
      if (c > 0x2300 && c < 0x23FF) {
        safeChars.add(c);
        continue;
      }

      // Miscellaneous Symbols
      if (c > 0x2600 && c < 0x26FF) {
        safeChars.add(c);
        continue;
      }

      // Dingbats
      if (c > 0x2700 && c < 0x27BF) {
        safeChars.add(c);
        continue;
      }

      // Miscellaneous Symbols and Arrows
      if (c > 0x2B00 && c < 0x2BFF) {
        safeChars.add(c);
        continue;
      }

      continue;
    }

    if ((c >= _a && c <= _z) || (c >= _A && c <= _Z) || (c >= _0 && c <= _9)) {
      safeChars.add(c);
      continue;
    }

    // Check for the other space characters, but add ' ' if found instaid of the spacing character
    if (spaceChars.contains(c)) {
      if (safeChars.isNotEmpty && !spaceChars.contains(safeChars.last)) {
        safeChars.add(c);
      }
      continue;
    }

    // Check for the other space characters, but add ' ' if found instaid of the spacing character
    if (otherSpaceChars.contains(c)) {
      if (safeChars.isNotEmpty && !spaceChars.contains(safeChars.last)) {
        safeChars.add(spaceChars.first);
      }
      continue;
    }

    safeChars.add(c);
  }
  return String.fromCharCodes(safeChars);
}
