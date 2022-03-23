import 'dart:async';
import 'dart:io';
import 'dart:convert';

import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';
import 'package:image/image.dart' as image;
import 'language_widgets.dart';
import 'info_widgets.dart';
import 'cv.dart';
import 'layout.dart';
import 'header.dart';
import 'footer.dart';
import 'args.dart';
import 'fonts.dart';

Future<void> main(List<String> programArgs) async {
  final ArgsParser args = ArgsParser(programArgs);
  final Style style = args.style;

  final CV cv;
  if (args.dummy) {
    print("using dummy data to create pdf");
    cv = CV.example();
  } else if (args.data.length != 0) {
    print("using data provided by argument to create pdf");
    final cvJsonData = jsonDecode(args.data);
    cv = CV.fromJson(cvJsonData);
  } else {
    print("did not provide the --data nor --dummy flag");
    exit(1);
  }

  // We need custom fonts as the default fon't doesn't have a lot of glyphs (sepcial characters)
  // The pdf library panics if a glyph is missing
  // As we handle with scraped data it's very common to see wired glyphs so if we want to create pdfs for those we'll need to use a custom font
  FontsManager fonts = FontsManager(args);

  List<dynamic> awaitedResults = await Future.wait([
    obtainLogo(args.logoImageUrl),
    fonts.resolvedFontRegular,
    fonts.resolvedFontBold,
  ]);

  image.Image? logo = awaitedResults[0];
  Font baseFont = awaitedResults[1];
  Font boldFont = awaitedResults[2];

  final pdf = Document(
    title: "CV",
    theme: ThemeData.withFont(
      base: baseFont,
      bold: boldFont,

      // Use the google icons font as the icons font
      icons: await fonts.iconsFont,
    ),
  );

  final List<ListWithHeader> largeListLayout = [];

  if (cv.workExperiences != null && cv.workExperiences!.isNotEmpty) {
    largeListLayout.add(ListWithHeader(
      IconData(0xe943), // Work
      "Werkervaring",
      cv.workExperiences!.map((workExper) => WorkExpWidget(workExper)).toList(),
    ));
  }

  if (cv.educations != null && cv.educations!.isNotEmpty) {
    largeListLayout.add(ListWithHeader(
      IconData(0xe80c), // School
      "Opleidingen",
      cv.educations!.map((education) => EducationWidget(education)).toList(),
    ));
  }
  if (cv.courses != null && cv.courses!.isNotEmpty) {
    largeListLayout.add(ListWithHeader(
      IconData(0xe865), // Book
      "Cursussen",
      cv.courses!.map((course) => EducationWidget(course)).toList(),
    ));
  }

  if (cv.languages != null && cv.languages!.isNotEmpty) {
    // The language list a very short widget so we can always add it to the remainingLists
    // (The remainingLists only shows small lists)
    largeListLayout.add(
      ListWithHeader(
        IconData(0xe8e2), // Translate
        "Talen",
        [
          LanguageLevelInfoWidget(style),
          ...cv.languages!.map((lang) => LanguageWidget(lang)).toList()
        ],
        canFitOnPage: true,
      ),
    );
  }

  final List<Widget> largeListLayoutWidgets = largeListLayout
      .map((list) => ColumnLayoutBlock(
            list,
            style,
            layoutBlockBasePadding.copyWith(top: PdfPageFormat.cm),
          ))
      .toList();

  pdf.addPage(
    MultiPage(
      footer: (Context context) => FooterWidget(
        ref: cv.referenceNumber,
        logo: logo,
        companyName: args.companyName,
        companyAddress: args.companyAddress,
      ),
      margin: const EdgeInsets.only(bottom: PdfPageFormat.cm),
      build: (Context context) => [
        HeaderWidget(cv: cv, style: style),
        Presentation(presentation: cv.presentation),
        ClientInfo(
          personalInfo: cv.personalDetails,
          driversLicenses: cv.driversLicenses,
        ),
        ...largeListLayoutWidgets,
      ],
    ),
  );

  final file = File(args.out);
  await file.writeAsBytes(await pdf.save());
}
