import 'dart:async';
import 'dart:io';
import 'dart:convert';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';
import 'package:args/args.dart';

import 'language_widgets.dart';
import 'info_widgets.dart';
import 'cv.dart';
import 'layout.dart';
import 'utils.dart';
import 'header.dart';
import 'footer.dart';

Future<void> main(List<String> args) async {
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

  var logo = await obtainLogo(argResult['logo-image-url']);
  BgColor headerColor = BgColor(PdfColor.fromHex(argResult['header-color']));
  BgColor subHeaderColor =
      BgColor(PdfColor.fromHex(argResult['sub-header-color']));

  final pdf = Document(
    title: "CV",
    theme: ThemeData.withFont(
      // We need custom fonts as the default fon't doesn't have a lot of glyphs (sepcial characters)
      // The pdf library panics if a glyph is missing
      // As we handle with scraped data it's very common to see wired glyphs so if we want to create pdfs for those we'll need to use a custom font
      base: await loadFont("./fonts/OpenSans-Regular.ttf"),
      bold: await loadFont("./fonts/OpenSans-Bold.ttf"),

      // Use the google icons font as the icons font
      icons: await loadFont("./fonts/MaterialIcons-Regular.ttf"),
    ),
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
      wrapLayoutBlocks.add(WrapLayoutBlock(list, subHeaderColor));
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
      footer: (Context context) => FooterWidget(
        ref: cv.referenceNumber,
        logo: logo,
        companyName: argResult['company-name'],
        companyAddress: argResult['company-address'],
      ),
      margin: const EdgeInsets.only(bottom: PdfPageFormat.cm),
      build: (Context context) => [
        HeaderWidget(cv: cv, headerColor: headerColor),
        LayoutBlockBase(
          child: ClientInfo(
            personalInformation: cv.personalDetails,
            driversLicenses: cv.driversLicenses,
          ),
        ),
        ColumnsLayoutBlock(remainingLists, subHeaderColor),
        ...wrapLayoutBlocks,
      ],
    ),
  );

  final file = File(argResult['out']);
  await file.writeAsBytes(await pdf.save());
}
