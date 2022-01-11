import 'package:pdf/pdf.dart';

String? formatDate(DateTime? input) {
  if (input == null) return null;

  String year = input.year.toString();
  String month = input.month.toString().padLeft(2, '0');
  String day = input.day.toString().padLeft(2, '0');

  return "${year}-${month}-${day}";
}

String? formatDateTime(DateTime? input) {
  if (input == null) return null;

  String hour = input.hour.toString().padLeft(2, '0');
  String minute = input.minute.toString().padLeft(2, '0');
  String second = input.second.toString().padLeft(2, '0');

  return "${formatDate(input)!} ${hour}:${minute}:${second}";
}

PdfColor getTextColorFromBg(PdfColor bgColor) {
  double red = bgColor.red * 255;
  double green = bgColor.green * 255;
  double blue = bgColor.blue * 255;

  // Calculate the luminance of the background color and base the text color on that
  // Colors are a shitshow and looking tough some Q and A's on the internet this calculation might be wrong
  // For a way to long thread about this issue see: https://stackoverflow.com/questions/596216/formula-to-determine-perceived-brightness-of-rgb-color
  double luminance = (red * 0.3) + (green * 0.59) + (blue * 0.11);
  return (luminance >= 128) ? PdfColors.black : PdfColors.white;
}

String? guessPostalCodePlace(String postalCode) {
  if (postalCode.length < 4) return null;

  int nr = int.parse(postalCode.substring(0, 4));

  if (nr < 1000) {
    return null;
  } else if (nr < 2000) {
    if (nr > 1900) return 'Castricum';
    if (nr > 1800) return 'Alkmaar';
    if (nr > 1700) return 'Heerhugowaard';
    if (nr > 1600) return 'Enkhuizen, Hoorn';
    if (nr > 1500) return 'Zaandam';
    if (nr > 1400) return 'Bussum, Uithoorn, Purmerend';
    if (nr > 1300) return 'Almere';
    if (nr > 1200) return 'Hilversum';
    if (nr > 1100) return 'Amsterdam, Amstelveen';
    if (nr > 1000) return 'Amsterdam';
  } else if (nr < 3000) {
    if (nr > 2900) return 'Capelle aan den IJssel';
    if (nr > 2800) return 'Gouda';
    if (nr > 2700) return 'Zoetermeer';
    if (nr > 2600) return 'Delft';
    if (nr > 2500) return 'Den Haag';
    if (nr > 2400) return 'Alphen aan den Rijn';
    if (nr > 2300) return 'Leiden';
    if (nr > 2200) return 'Katwijk';
    if (nr > 2100) return 'Heemstede';
    if (nr > 2000) return 'Haarlem';
  } else if (nr < 4000) {
    if (nr > 3900) return 'Veenendaal';
    if (nr > 3800) return 'Amersfoort';
    if (nr > 3700) return 'Zeist';
    if (nr > 3600) return 'Maarssen';
    if (nr > 3500) return 'Utrecht (stad)';
    if (nr > 3400) return 'IJsselstein';
    if (nr > 3300) return 'Dordrecht, Drechtsteden';
    if (nr > 3200) return 'Spijkenisse';
    if (nr > 3100) return 'Schiedam, Vlaardingen';
    if (nr > 3000) return 'Rotterdam';
  } else if (nr < 5000) {
    if (nr > 4900) return 'Oosterhout';
    if (nr > 4800) return 'Breda';
    if (nr > 4700) return 'Roosendaal';
    if (nr > 4600) return 'Bergen op Zoom';
    if (nr > 4500) return 'Oostburg';
    if (nr > 4400) return 'Yerseke';
    if (nr > 4300) return 'Zierikzee';
    if (nr > 4200) return 'Gorinchem';
    if (nr > 4100) return 'Culemborg';
    if (nr > 4000) return 'Tiel';
  } else if (nr < 6000) {
    if (nr > 5900) return 'Venlo';
    if (nr > 5800) return 'Venray';
    if (nr > 5700) return 'Helmond';
    if (nr > 5600) return 'Eindhoven';
    if (nr > 5500) return 'Veldhoven';
    if (nr > 5400) return 'Uden';
    if (nr > 5300) return 'Zaltbommel';
    if (nr > 5200) return '\'s-Hertogenbosch';
    if (nr > 5100) return 'Dongen';
    if (nr > 5000) return 'Tilburg';
  } else if (nr < 7000) {
    if (nr > 6900) return 'Zevenaar';
    if (nr > 6800) return 'Arnhem';
    if (nr > 6700) return 'Wageningen, Ede';
    if (nr > 6600) return 'Wijchen';
    if (nr > 6500) return 'Nijmegen';
    if (nr > 6400) return 'Heerlen';
    if (nr > 6300) return 'Valkenburg aan de Geul';
    if (nr > 6200) return 'Maastricht';
    if (nr > 6100) return 'Echt, Limburg';
    if (nr > 6000) return 'Weert';
  } else if (nr < 8000) {
    if (nr > 7900) return 'Hoogeveen';
    if (nr > 7800) return 'Emmen';
    if (nr > 7700) return 'Dedemsvaart';
    if (nr > 7600) return 'Almelo';
    if (nr > 7500) return 'Enschede';
    if (nr > 7400) return 'Deventer';
    if (nr > 7300) return 'Apeldoorn';
    if (nr > 7200) return 'Zutphen';
    if (nr > 7100) return 'Winterswijk';
    if (nr > 7000) return 'Doetinchem';
  } else if (nr < 9000) {
    if (nr > 8900) return 'Leeuwarden';
    if (nr > 8800) return 'Franeker';
    if (nr > 8700) return 'Bolsward';
    if (nr > 8600) return 'Sneek';
    if (nr > 8500) return 'Joure';
    if (nr > 8400) return 'Gorredijk';
    if (nr > 8300) return 'Emmeloord';
    if (nr > 8200) return 'Lelystad';
    if (nr > 8100) return 'Raalte';
    if (nr > 8000) return 'Zwolle';
  } else {
    if (nr > 9900) return 'Appingedam';
    if (nr > 9800) return 'Zuidhorn';
    if (nr > 9700) return 'Groningen (stad)';
    if (nr > 9600) return 'Hoogezand-Sappemeer';
    if (nr > 9500) return 'Stadskanaal';
    if (nr > 9400) return 'Assen';
    if (nr > 9300) return 'Roden';
    if (nr > 9200) return 'Drachten';
    if (nr > 9100) return 'Dokkum';
    if (nr > 9000) return 'Grouw (Grou)';
  }
}
