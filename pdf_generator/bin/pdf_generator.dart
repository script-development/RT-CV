import 'dart:async';
import 'dart:io';
import 'dart:convert';
import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'language_widgets.dart';
import 'info_widgets.dart';
import 'cv.dart';
import 'layout.dart';
import 'utils.dart';
import 'header.dart';
import 'footer.dart';
import 'args.dart';
import 'fonts.dart';

Future<void> main(List<String> programArgs) async {
  final ArgsParser args = ArgsParser(programArgs);

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

  var logo = await obtainLogo(args.logoImageUrl);
  final BgColor headerColor = BgColor(args.headerColor);
  final BgColor subHeaderColor = BgColor(args.subHeaderColor);

  // We need custom fonts as the default fon't doesn't have a lot of glyphs (sepcial characters)
  // The pdf library panics if a glyph is missing
  // As we handle with scraped data it's very common to see wired glyphs so if we want to create pdfs for those we'll need to use a custom font
  FontsManager fonts = FontsManager(args);

  final pdf = Document(
    title: "CV",
    theme: ThemeData.withFont(
      base: await fonts.resolvedFontRegular,
      bold: await fonts.resolvedFontBold,

      // Use the google icons font as the icons font
      icons: await fonts.iconsFont,
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
        companyName: args.companyName,
        companyAddress: args.companyAddress,
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

  final file = File(args.out);
  await file.writeAsBytes(await pdf.save());
}
