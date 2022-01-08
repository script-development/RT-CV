import 'dart:async';
import 'dart:io';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';
import 'package:args/args.dart';
import 'dart:typed_data';
import 'dart:convert';

import 'language_widgets.dart';
import 'cv.dart';
import 'layout.dart';
import 'utils.dart';

Future<void> main(List<String> args) async {
  ArgParser argsParser = ArgParser();
  argsParser.addFlag(
    'help',
    abbr: 'h',
    help: "Print this message",
  );
  argsParser.addFlag(
    'dummy',
    abbr: 'D',
    help:
        "Use dummy data, handy for working on this application. The dummy data is located in bin/cv.dart",
  );
  argsParser.addOption(
    'data',
    abbr: 'd',
    help:
        'input CV as json data (the structure of this data should be the json of the CV structure in ../models/cv.go)',
  );
  argsParser.addOption(
    'out',
    abbr: 'o',
    defaultsTo: 'example.pdf',
    help: "to where should we write the output file",
  );

  final ArgResults argResult = argsParser.parse(args);
  if (argResult['help'] == true) {
    print(argsParser.usage);
    exit(0);
  }

  final String dataFlag = argResult['data'] ?? '';

  final CV cv;
  if (argResult['dummy']) {
    print("using dummy data to create pdf");
    cv = CV.example();
  } else if (dataFlag.length != 0) {
    print("using data provided by argument to create pdf");
    final cvJsonData = jsonDecode(dataFlag);
    cv = CV.fromJson(cvJsonData);
  } else {
    print("did not provide the --data nor --dummy flag");
    exit(1);
  }

  final File fontFile = File("./MaterialIcons-Regular.ttf");
  Uint8List data = await fontFile.readAsBytesSync();
  Font materialIconsFont = Font.ttf(ByteData.view(data.buffer));

  final pdf = Document(
    title: "CV",
    theme: ThemeData.withFont(icons: materialIconsFont),
  );

  final List<ListWithHeader> lists = [];

  if (cv.workExperiences != null && cv.workExperiences!.isNotEmpty) {
    lists.add(ListWithHeader(
      IconData(0xe943), // Work
      "Werkervaring",
      cv.workExperiences!.map((workExper) => WorkExpWidget(workExper)).toList(),
    ));
  }
  if (cv.educations != null && cv.educations!.isNotEmpty) {
    lists.add(ListWithHeader(
      IconData(0xe80c), // School
      "Opleidingen",
      cv.educations!.map((education) => EducationWidget(education)).toList(),
    ));
  }
  if (cv.courses != null && cv.courses!.isNotEmpty) {
    lists.add(ListWithHeader(
      IconData(0xe865), // Book
      "Cursussen",
      cv.courses!.map((course) => CourseWidget(course)).toList(),
    ));
  }

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

  if (cv.languages != null && cv.languages!.isNotEmpty) {
    // The language list a very short widget so we can always add it to the remainingLists
    // (The remainingLists only shows small lists)
    remainingLists.add(
      ListWithHeader(
        IconData(0xe8e2), // Translate
        "Talen",
        [
          LanguageLevelInfoWidget(),
          ...cv.languages!.map((lang) => LanguageWidget(lang)).toList()
        ],
      ),
    );
  }

  pdf.addPage(
    MultiPage(
      margin: const EdgeInsets.only(bottom: PdfPageFormat.cm),
      build: (Context context) => [
        Header(cv),
        LayoutBlockBase(
          child: ClientInfo(
            personalInformation: cv.personalDetails,
            driversLicenses: cv.driversLicenses,
          ),
        ),
        ColumnsLayoutBlock(remainingLists),
        ...wrapLayoutBlocks,
      ],
    ),
  );

  final file = File(argResult['out']);
  await file.writeAsBytes(await pdf.save());
}

class Header extends StatelessWidget {
  Header(this.cv);

  final CV cv;

  @override
  Widget build(Context context) {
    var metaColor = TextStyle(
      fontSize: 6,
      color: PdfColor.fromInt(0xffb3d8dd),
    );

    List<Widget> meta = [
      Text("ref #" + cv.referenceNumber, style: metaColor),
    ];
    if (cv.lastChanged != null) {
      meta.add(Text(
        "laatst geupdate " + formatDateTime(cv.lastChanged)!,
        style: metaColor,
      ));
    } else if (cv.createdAt != null) {
      meta.add(Text(
        "cv gemaakt op " + formatDateTime(cv.createdAt)!,
        style: metaColor,
      ));
    }

    return SizedBox(
      width: double.infinity,
      child: Container(
        color: PdfColor.fromInt(0xff4ca1af),
        padding: const EdgeInsets.symmetric(
          horizontal: PdfPageFormat.cm,
          vertical: PdfPageFormat.cm * 1.5,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              cv.personalDetails.fullName,
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: PdfColors.white,
              ),
            ),
            ...meta,
          ],
        ),
      ),
    );
  }
}

class ClientInfo extends StatelessWidget {
  ClientInfo({
    required this.personalInformation,
    this.driversLicenses,
  });

  List<Widget> children = [];
  final PersonalDetails personalInformation;
  final List<String>? driversLicenses;

  final TextStyle labelStyle = TextStyle(
    fontSize: 8,
    color: PdfColors.grey,
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
    tryAddToList("Email", personalInformation.email);
    tryAddToList("Telefoon", personalInformation.phoneNumber);
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

    if (!personalInformation.hasAddress) {
      tryAddToList("Postcode", personalInformation.zip);
      return Wrap(children: children, spacing: 10);
    }

    return Row(
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
                    Text("Stad: ", style: labelStyle),
                    Text(personalInformation.city!, style: valueStyle),
                  ],
                ),
                Row(
                  children: [
                    Text("Address: ", style: labelStyle),
                    Text(
                        "${personalInformation.streetName} ${personalInformation.houseNumber} ${personalInformation.houseNumberSuffix}",
                        style: valueStyle),
                  ],
                ),
                Row(
                  children: [
                    Text("Postcode: ", style: labelStyle),
                    Text(personalInformation.zip!, style: valueStyle),
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
    );
  }
}

class WorkExpWidget extends StatelessWidget {
  WorkExpWidget(WorkExperience this.exp);

  final WorkExperience exp;

  @override
  Widget build(Context context) {
    return ListEntry(
      exp.profession ?? '??',
      company: exp.employer,
      description: exp.description,
      from: exp.startDate,
      to: exp.endDate,
    );
  }
}

class CourseWidget extends StatelessWidget {
  CourseWidget(Course this.course);

  final Course course;

  @override
  Widget build(Context context) {
    return ListEntry(
      course.name,
      company: course.institute,
      from: course.startDate,
      to: course.endDate,
      description: course.description,
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
      company: education.institute,
      from: education.startDate,
      to: education.endDate,
      description: education.description,
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
      padding: const EdgeInsets.only(top: 5),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: children,
      ),
    );
  }
}
