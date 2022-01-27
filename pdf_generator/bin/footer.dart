import 'dart:convert';
import 'package:pdf/pdf.dart';
import 'dart:io';
import 'package:pdf/widgets.dart';
import 'package:image/image.dart' as image;
import 'package:http/http.dart' as http;

import 'utils.dart';

class FooterWidget extends StatelessWidget {
  FooterWidget({
    required this.ref,
    this.companyName,
    this.companyAddress,
    this.logo,
  });

  final String ref;
  final String? companyName;
  final String? companyAddress;
  final image.Image? logo;

  final TextStyle footerLabelStyle = TextStyle(
    fontSize: 7,
    color: PdfColors.grey500,
  );
  final TextStyle footerValueStyle = TextStyle(
    fontSize: 7,
    color: PdfColors.grey700,
  );

  Widget labelAndValue(String label, String value) {
    return Row(
      children: [
        Text(label + " ", style: footerLabelStyle),
        Text(value, style: footerValueStyle),
      ],
    );
  }

  @override
  Widget build(Context context) {
    List<Widget> rowWidgets = [
      Expanded(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            labelAndValue(
              "Pagina",
              "${context.pageNumber} van ${context.pagesCount}",
            ),
            labelAndValue(
              "PDF Gecreëerd op",
              formatDateTime(DateTime.now()) ?? '--',
            ),
            labelAndValue(
              "Ref",
              ref,
            ),
          ],
        ),
      ),
    ];

    if (logo != null || companyName != null || companyAddress != null) {
      bool placeNameAboveLogo = companyName != null && companyAddress == null;
      bool placeNameAboveAddress =
          companyName != null && companyAddress != null;

      if (logo != null) {
        rowWidgets.add(Expanded(
          child: Column(
            crossAxisAlignment: companyAddress != null
                ? CrossAxisAlignment.center
                : CrossAxisAlignment.end,
            children: [
              ...(placeNameAboveLogo
                  ? [
                      Padding(
                        padding: EdgeInsets.only(bottom: 3),
                        child: Text(companyName!, style: footerValueStyle),
                      )
                    ]
                  : []),
              Image(
                ImageImage(logo!),
                fit: BoxFit.scaleDown,
                height: 20,
              ),
            ],
          ),
        ));
      }

      if (companyAddress != null) {
        rowWidgets.add(Expanded(
          child: Container(
            alignment: Alignment.centerRight,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                ...(placeNameAboveAddress
                    ? [
                        Padding(
                          padding: EdgeInsets.only(bottom: 3),
                          child: Text(companyName!, style: footerValueStyle),
                        )
                      ]
                    : []),
                Text(
                  "Adres",
                  style: footerLabelStyle,
                ),
                Text(
                  companyAddress!,
                  style: footerValueStyle,
                  textAlign: TextAlign.right,
                  overflow: TextOverflow.clip,
                ),
              ],
            ),
          ),
        ));
      }

      if (companyName != null && logo == null && companyAddress == null) {
        rowWidgets.add(Expanded(
          child: Container(
            alignment: Alignment.centerRight,
            child: Text(companyName!, style: footerValueStyle),
          ),
        ));
      }
    }

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: PdfPageFormat.cm),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: rowWidgets,
      ),
    );
  }
}

Future<image.Image?> obtainLogo(String? url) async {
  if (url == null ||
      (!url.startsWith("http://") && !url.startsWith("https://"))) {
    return null;
  }

  List<int> urlHash = List.generate(32, (index) => 0);
  var urlChars = utf8.encode(url);
  for (var i = 0; i < urlChars.length; i++) {
    var urlHashIdx = i % urlHash.length;
    urlHash[urlHashIdx] ^= urlChars[i];
  }

  var logoCacheFile = File(
    'cache/' +
        urlHash.map((e) => e.toRadixString(16)).join('').padLeft(2, '0') +
        '.png',
  );

  try {
    if (await logoCacheFile.exists()) {
      var logoBytes = await logoCacheFile.readAsBytes();
      return image.decodePng(logoBytes);
    }

    http.Response response = await http.get(Uri.parse(url));
    var recivedImage = image.decodeImage(response.bodyBytes);
    if (recivedImage == null) return null;

    if (recivedImage.height > 100) {
      recivedImage = image.copyResize(
        recivedImage,
        height: 50,
        interpolation: image
            .Interpolation.average, // average and linear seem to work best here
      );
    }

    var cacheDir = Directory("cache");
    if (!await cacheDir.exists()) {
      await cacheDir.create();
    }

    await logoCacheFile.writeAsBytes(image.encodePng(recivedImage));
    return recivedImage;
  } catch (e) {
    print("An error occurred while loading the logo: ${e}");
    print("Continuing without logo");
    return null;
  }
}
