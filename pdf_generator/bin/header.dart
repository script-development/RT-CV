import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'cv.dart';
import 'utils.dart';

class HeaderWidget extends StatelessWidget {
  HeaderWidget({
    required this.cv,
    required this.headerColor,
  });

  final CV cv;
  final BgColor headerColor;

  @override
  Widget build(Context context) {
    final PdfColor textColor = headerColor.textColor;
    final PdfColor bgColor = headerColor.bgColor;

    // Create meta color
    PdfColorHsl bgColorAsHsl = headerColor.bgColor.toHsl();
    final double hue = bgColorAsHsl.hue;
    final double saturation = bgColorAsHsl.saturation;
    double lightness = bgColorAsHsl.lightness;
    if (textColor == PdfColors.white) {
      lightness += .35;
      if (lightness > 1) {
        lightness = 1;
      }
    } else {
      lightness -= .35;
      if (lightness < 0) {
        lightness = 0;
      }
    }

    var metaColor = TextStyle(
      fontSize: 6,
      color: PdfColorHsl(hue, saturation, lightness),
    );

    List<Widget> meta = [];
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
        color: bgColor,
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
                color: textColor,
              ),
            ),
            ...meta,
          ],
        ),
      ),
    );
  }
}
