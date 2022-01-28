import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'cv.dart';
import 'utils.dart';
import 'args.dart';

class HeaderWidget extends StatelessWidget {
  HeaderWidget({
    required this.cv,
    required this.style,
  });

  final CV cv;
  final Style style;

  @override
  Widget build(Context context) {
    // Create meta color
    PdfColorHsl bgColorAsHsl = style.headerBackgroundColor.toHsl();
    final double hue = bgColorAsHsl.hue;
    double saturation = bgColorAsHsl.saturation;
    double lightness = bgColorAsHsl.lightness;
    if (style.headerTextColor == PdfColors.white) {
      if (lightness < .1 && hue < .1) {
        // Fix color of text becomming red when a #000 (black) is provided as color
        // This is a side effect of converting hex to hsl that is fixed here.
        saturation = 0;
      }
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
        color: style.headerBackgroundColor,
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
                color: style.headerTextColor,
              ),
            ),
            ...meta,
          ],
        ),
      ),
    );
  }
}
