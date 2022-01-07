import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'data.dart';

const readingColor = PdfColors.green400;
const writingColor = PdfColors.blue400;

class LanguageLevelInfoWidget extends StatelessWidget {
  final TextStyle labelStyle = TextStyle(
    fontSize: 8,
    color: PdfColors.grey700,
  );

  Widget labelLanguageSkill(LanguageSkillLevel languageSkillLevel) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        Text(
          humanLanguageSkillLevel[languageSkillLevel] ?? '',
          style: labelStyle,
        ),
        Container(
          color: PdfColors.grey700,
          height: 5,
          width: 1,
        ),
      ],
    );
  }

  Widget labelSkillKind(String kind, PdfColor color) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          kind + " ",
          style: labelStyle,
        ),
        Container(
          decoration: BoxDecoration(
            color: color,
            borderRadius: BorderRadius.circular(4),
          ),
          height: 8,
          width: 8,
        ),
      ],
    );
  }

  @override
  Widget build(Context context) {
    return Column(
      children: [
        Container(
          margin: const EdgeInsets.only(top: 5),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                padding: const EdgeInsets.only(right: 5),
                child: labelSkillKind("lezen", readingColor),
              ),
              labelSkillKind("schijven", writingColor),
            ],
          ),
        ),
        Container(
          margin: const EdgeInsets.only(top: 5),
          child: Row(
            children: [
              Container(
                constraints: BoxConstraints.tightFor(width: 70),
                child: labelLanguageSkill(LanguageSkillLevel.unknown),
              ),
              Expanded(
                child: labelLanguageSkill(LanguageSkillLevel.reasonable),
              ),
              Expanded(
                child: labelLanguageSkill(LanguageSkillLevel.good),
              ),
              Expanded(
                child: labelLanguageSkill(LanguageSkillLevel.excellent),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class LanguageSkillWidget extends StatelessWidget {
  LanguageSkillWidget(this.languageSkill);

  final LanguageSkill languageSkill;

  int get readingNr => languageSkillLevelNr[languageSkill.reading] ?? 0;
  int get writingNr => languageSkillLevelNr[languageSkill.writing] ?? 0;

  @override
  Widget build(Context context) {
    return Container(
      margin: EdgeInsets.only(top: 4),
      child: Row(
        children: [
          Container(
            constraints: BoxConstraints.tightFor(width: 70),
            child: Text(
              languageSkill.name,
              overflow: TextOverflow.clip,
              style: TextStyle(
                fontSize: 10,
              ),
            ),
          ),
          Expanded(
            child: Container(
              height: 10,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(5),
                color: PdfColors.grey200,
              ),
              child: Stack(
                children: [
                  LanguageLevelbar(readingNr, readingColor, 10),
                  LanguageLevelbar(writingNr, writingColor, 10),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class LanguageLevelbar extends StatelessWidget {
  LanguageLevelbar(this.levelNr, this.color, this.height);

  final int levelNr;
  final PdfColor color;
  final double height;

  @override
  Widget build(Context context) {
    if (levelNr == 0) {
      // Display a simple dot to indicate that the language is unknown.
      return Container(
        height: height,
        width: height,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(5),
          color: color,
        ),
      );
    } else if (levelNr == maxLanguageLevelNr) {
      return Expanded(
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(5),
            color: color,
          ),
        ),
      );
    } else {
      return Row(
        children: [
          Expanded(
            flex: levelNr,
            child: Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(5),
                color: color,
              ),
            ),
          ),
          Expanded(
            flex: maxLanguageLevelNr - levelNr,
            child: Container(),
          ),
        ],
      );
    }
  }
}
