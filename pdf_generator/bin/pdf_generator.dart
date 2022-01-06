import 'dart:async';
import 'dart:io';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';
import 'dart:typed_data';

import './data.dart';
import 'layout.dart';

Future<void> main(List<String> arguments) async {
  var fontFile = File("./MaterialIcons-Regular.ttf");
  Uint8List data = await fontFile.readAsBytesSync();
  Font materialIconsFont = Font.ttf(ByteData.view(data.buffer));

  final pdf = Document(
    title: "CV",
    theme: ThemeData.withFont(icons: materialIconsFont),
  );

  CV cv = CV.example();

  List<ListWithHeader> lists = [
    ListWithHeader(
      "Werkervaring",
      cv.workExpr.map((workExpr) => WorkExpWidget(workExpr)).toList(),
    ),
    ListWithHeader(
      "Opleidingen",
      cv.education.map((education) => EducationWidget(education)).toList(),
    ),
    ListWithHeader(
      "Cursussen",
      cv.courses.map((education) => EducationWidget(education)).toList(),
    ),
    ListWithHeader(
      "Talen",
      cv.languageSkills.map((skill) => LanguageSkillWidget(skill)).toList(),
    ),
  ];

  // Determain the layout depending on the amound of items in the lists.
  List<WrapLayoutBlock> wrapLayoutBlocks = [];
  List<ListWithHeader> remainingLists = [];
  for (ListWithHeader list in lists) {
    if (list.length > 4) {
      wrapLayoutBlocks.add(WrapLayoutBlock(list));
    } else {
      remainingLists.add(list);
    }
  }

  pdf.addPage(
    MultiPage(
      margin: EdgeInsets.only(bottom: PdfPageFormat.cm),
      build: (Context context) => [
        Header(
          name: cv.name,
          ref: cv.reference,
        ),
        LayoutBlockBase(
          child: ClientInfo(
            email: cv.email,
            phoneNr: cv.phoneNr,
            driversLicenses: cv.driversLicenses,
          ),
        ),
        ColumnsLayoutBlock(remainingLists),
        ...wrapLayoutBlocks,
      ],
    ),
  );

  final file = File('example.pdf');
  await file.writeAsBytes(await pdf.save());
}

class Header extends StatelessWidget {
  Header({
    required this.name,
    required this.ref,
  });

  final String name;
  final String ref;

  @override
  Widget build(Context context) {
    return SizedBox(
      width: double.infinity,
      child: Container(
        color: PdfColor.fromInt(0xff4ca1af),
        padding: EdgeInsets.symmetric(
          horizontal: PdfPageFormat.cm,
          vertical: PdfPageFormat.cm * 1.5,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              name,
              overflow: TextOverflow.clip,
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: PdfColors.white,
              ),
            ),
            Text(
              "ref #" + ref,
              overflow: TextOverflow.clip,
              style: TextStyle(
                fontSize: 6,
                color: PdfColor.fromInt(0xffb3d8dd),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ClientInfo extends StatelessWidget {
  ClientInfo({
    this.email,
    this.phoneNr,
    this.driversLicenses,
  });

  List<Widget> children = [];
  final String? email;
  final String? phoneNr;
  final List<String>? driversLicenses;

  tryAddToList(String label, String? value) {
    if (value != null) {
      children.add(
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              label + ": ",
              overflow: TextOverflow.clip,
              style: TextStyle(
                fontSize: 8,
                color: PdfColors.grey,
              ),
            ),
            Text(
              value,
              overflow: TextOverflow.clip,
              style: TextStyle(
                fontSize: 10,
                color: PdfColors.black,
              ),
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(Context context) {
    children = [];
    tryAddToList("Email", email);
    tryAddToList("Telefoon", phoneNr);
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

    return Wrap(children: children, spacing: 10);
  }
}

class LanguageSkillWidget extends StatelessWidget {
  LanguageSkillWidget(this.languageSkill);

  final LanguageSkill languageSkill;

  get humandReading => humanLanguageSkillLevel[languageSkill.reading];
  get humandWriting => humanLanguageSkillLevel[languageSkill.writing];

  @override
  Widget build(Context context) {
    return ListEntry(
      languageSkill.name,
      description: "Lezen: ${humandReading}, Schijven: ${humandWriting}",
    );
  }
}

class WorkExpWidget extends StatelessWidget {
  WorkExpWidget(WorkExp this.workExp);

  final WorkExp workExp;

  @override
  Widget build(Context context) {
    return ListEntry(
      workExp.name,
      company: workExp.company,
      description: workExp.description,
      from: workExp.from,
      to: workExp.to,
    );
  }
}

class EducationWidget extends StatelessWidget {
  EducationWidget(Education this.education);

  final Education education;

  @override
  Widget build(Context context) {
    return ListEntry(
      education.name,
      company: education.org,
      from: education.from,
      to: education.to,
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
  });

  final String title;
  final String? company;
  final String? description;
  final DateTime? from;
  final DateTime? to;

  String? formatDate(DateTime? input) {
    if (input == null) return null;

    String year = input.year.toString();
    String month = input.month.toString().padLeft(2, '0');
    String day = input.day.toString().padLeft(2, '0');

    return "${year}-${month}-${day}";
  }

  @override
  Widget build(Context context) {
    List<Widget> children = [
      Flexible(
        child: Text(
          title,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontSize: 10,
          ),
        ),
      ),
    ];

    TextStyle contentStyle = TextStyle(
      fontSize: 8,
      color: PdfColors.grey800,
    );
    TextStyle labelStyle = TextStyle(
      fontSize: 8,
      color: PdfColors.grey600,
    );

    if (company != null && company!.isNotEmpty) {
      children.add(
        Row(children: [
          Text("Bij: ", style: labelStyle),
          Text(
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
      children.add(
        Row(children: [
          Text("Vanaf ", style: labelStyle),
          Text(fromStr ?? '??', style: contentStyle),
          Text(" tot ", style: labelStyle),
          Text(toStr ?? '??', style: contentStyle),
        ]),
      );
    }

    if (description != null && description!.isNotEmpty) {
      children.add(Flexible(
        child: Text(
          (description!.length > 300)
              ? description!.substring(0, 300) + '..'
              : description!,
          overflow: TextOverflow.clip,
          style: contentStyle,
        ),
      ));
    }

    return Padding(
      padding: EdgeInsets.only(top: 5),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: children,
      ),
    );
  }
}
