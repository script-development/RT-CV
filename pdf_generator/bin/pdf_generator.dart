import 'dart:async';
import 'dart:io';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import './data.dart';

Future<void> main(List<String> arguments) async {
  final pdf = Document();

  CV cv = CV.example();

  List<WorkExpWidget> workExpr =
      cv.workExpr.map((workExpr) => WorkExpWidget(workExpr)).toList();
  List<EducationWidget> educations =
      cv.education.map((education) => EducationWidget(education)).toList();
  List<EducationWidget> course =
      cv.courses.map((education) => EducationWidget(education)).toList();
  List<LanguageSkillWidget> languageSkills =
      cv.languageSkills.map((skill) => LanguageSkillWidget(skill)).toList();

  pdf.addPage(
    MultiPage(
      margin: EdgeInsets.all(0),
      build: (Context context) => [
        Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Container(
              color: PdfColor.fromInt(0xff4ca1af),
              padding: EdgeInsets.symmetric(
                horizontal: PdfPageFormat.cm,
                vertical: PdfPageFormat.cm * 1.5,
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    cv.name,
                    overflow: TextOverflow.clip,
                    style: TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                      color: PdfColors.white,
                    ),
                  ),
                  Text(
                    "ref #" + cv.reference,
                    overflow: TextOverflow.clip,
                    style: TextStyle(
                      fontSize: 6,
                      color: PdfColor.fromInt(0xffb3d8dd),
                    ),
                  ),
                ],
              ),
            ),
            Padding(
              padding: EdgeInsets.all(PdfPageFormat.cm),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  ClientInfo(
                      email: cv.email,
                      phoneNr: cv.phoneNr,
                      driversLicenses: cv.driversLicenses),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Padding(
                          padding: EdgeInsets.only(right: PdfPageFormat.cm / 2),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              ListTitle("Werkervaring"),
                              ...workExpr,
                              ListTitle("Talen"),
                              ...languageSkills,
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        child: Padding(
                          padding: EdgeInsets.only(left: PdfPageFormat.cm / 2),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              ListTitle("Opleidingen"),
                              ...educations,
                              ListTitle("Cursussen"),
                              ...course,
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ],
    ),
  );

  final file = File('example.pdf');
  await file.writeAsBytes(await pdf.save());
}

class ListTitle extends StatelessWidget {
  ListTitle(this.title);

  final String title;

  @override
  Widget build(Context context) {
    return Container(
      color: PdfColor.fromInt(0xffffe004),
      margin: EdgeInsets.only(top: 10),
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        child: Text(
          title,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontSize: 10,
            fontWeight: FontWeight.bold,
          ),
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
    return (input == null)
        ? null
        : "${input.year.toString()}-${input.month.toString().padLeft(2, '0')}-${input.day.toString().padLeft(2, '0')}";
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
          Text("At: ", style: labelStyle),
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
