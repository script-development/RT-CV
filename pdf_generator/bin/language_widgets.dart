import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'args.dart';
import 'cv.dart';

const _writingColor = PdfColors.green400;
const _speakingColor = PdfColors.blue400;

class LanguageLevelInfoWidget extends StatelessWidget {
  LanguageLevelInfoWidget(this.style);

  final Style style;

  final TextStyle labelStyle = TextStyle(
    fontSize: 8,
    color: PdfColors.grey700,
  );

  Widget labelLanguage(LanguageLevel languageLevel) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        Text(
          humanLanguageLevel[languageLevel] ?? '',
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

  Widget labelLevelKind(String kind, PdfColor color) {
    Widget dot = Container(
      decoration: BoxDecoration(
        color: color,
        borderRadius: BorderRadius.circular(4),
      ),
      height: 8,
      width: 8,
    );

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: style.layoutStyle == LayoutStyle.style_1
          ? [
              Text(
                kind + " ",
                style: labelStyle,
              ),
              dot,
            ]
          : [
              dot,
              Text(
                " " + kind,
                style: labelStyle,
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
            mainAxisAlignment: style.layoutStyle == LayoutStyle.style_1
                ? MainAxisAlignment.end
                : MainAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.only(right: 5),
                child: labelLevelKind("schrijven", _writingColor),
              ),
              labelLevelKind("spreken", _speakingColor),
            ],
          ),
        ),
        Container(
          margin: const EdgeInsets.only(top: 5),
          child: Row(
            children: [
              Container(
                constraints: BoxConstraints.tightFor(width: 70),
                child: labelLanguage(LanguageLevel.unknown),
              ),
              Expanded(
                child: labelLanguage(LanguageLevel.reasonable),
              ),
              Expanded(
                child: labelLanguage(LanguageLevel.good),
              ),
              Expanded(
                child: labelLanguage(LanguageLevel.excellent),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class LanguageWidget extends StatelessWidget {
  LanguageWidget(this.language);

  final Language language;

  int get writingNr => languageLevelToNr[language.levelWritten] ?? 0;
  int get speakingNr => languageLevelToNr[language.levelSpoken] ?? 0;

  @override
  Widget build(Context context) {
    return Container(
      margin: EdgeInsets.only(top: 4),
      child: Row(
        children: [
          Container(
            constraints: BoxConstraints.tightFor(width: 70),
            child: Text(
              language.name,
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
                children: speakingNr == writingNr
                    ? [
                        _LanguageLevelbar(
                          writingNr,
                          _writingColor,
                          10,
                          addOffsetToLeft: true,
                        ),
                        _LanguageLevelbar(
                          speakingNr,
                          _speakingColor,
                          10,
                          addOffsetToRight: true,
                        ),
                      ]
                    : speakingNr > writingNr
                        ? [
                            _LanguageLevelbar(speakingNr, _speakingColor, 10),
                            _LanguageLevelbar(writingNr, _writingColor, 10),
                          ]
                        : [
                            _LanguageLevelbar(writingNr, _writingColor, 10),
                            _LanguageLevelbar(speakingNr, _speakingColor, 10),
                          ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _LanguageLevelbar extends StatelessWidget {
  _LanguageLevelbar(
    this.levelNr,
    this.color,
    this.height, {
    this.addOffsetToRight,
    this.addOffsetToLeft,
  });

  final int levelNr;
  final PdfColor color;
  final double height;
  final bool? addOffsetToRight;
  final bool? addOffsetToLeft;

  double get offsetRight => addOffsetToRight == true ? 5 : 0;
  double get offsetLeft => addOffsetToLeft == true ? 5 : 0;

  @override
  Widget build(Context context) {
    var padding = EdgeInsets.only(
      left: offsetLeft,
      right: offsetRight,
    );

    if (levelNr == 0) {
      // Display a simple dot to indicate that the language is unknown.
      return Padding(
        padding: padding,
        child: Container(
          height: height,
          width: height,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(5),
            color: color,
          ),
        ),
      );
    } else if (levelNr == maxLanguageLevelNr) {
      return Padding(
        padding: padding,
        child: Expanded(
          child: Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(5),
              color: color,
            ),
          ),
        ),
      );
    } else {
      return Row(
        children: [
          Expanded(
            flex: levelNr,
            child: Padding(
              padding: padding,
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(5),
                  color: color,
                ),
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
