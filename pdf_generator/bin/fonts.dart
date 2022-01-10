import 'package:pdf/widgets.dart';
import 'dart:io';
import 'dart:typed_data';

import 'args.dart';

class _FontFiles {
  const _FontFiles(this.regular, this.bold);
  final String regular;
  final String bold;
}

class FontsManager extends ArgsParser {
  FontsManager(ArgsParser argsParser)
      : super.fromArgResults(argsParser.argResult);

  Future<Font> loadFont(String filename) async {
    final File fontFile = File("./fonts/" + filename);
    Uint8List data = await fontFile.readAsBytesSync();
    return Font.ttf(ByteData.view(data.buffer));
  }

  Future<Font> get iconsFont => loadFont("MaterialIcons-Regular.ttf");

  Future<Font> get resolvedFontRegular => loadFont(
        _getFontOrFallback(super.fontRegular).regular,
      );

  Future<Font> get resolvedFontBold => loadFont(
        _getFontOrFallback(super.fontBold).bold,
      );

  final Map<String, _FontFiles> _fontFilesMap = {
    'BeVietnamPro':
        _FontFiles("BeVietnamPro-Regular.ttf", "BeVietnamPro-Bold.ttf"),
    'IBMPlexMono':
        _FontFiles("IBMPlexMono-Regular.ttf", "IBMPlexMono-Bold.ttf"),
    'IBMPlexSans':
        _FontFiles("IBMPlexSans-Regular.ttf", "IBMPlexSans-Bold.ttf"),
    'IBMPlexSerif':
        _FontFiles("IBMPlexSerif-Regular.ttf", "IBMPlexSerif-Bold.ttf"),
    'Lobster': _FontFiles("Lobster-Regular.ttf", "Lobster-Regular.ttf"),
    'OpenSans': _FontFiles("OpenSans-Regular.ttf", "OpenSans-Bold.ttf"),
    'PlayfairDisplay':
        _FontFiles("PlayfairDisplay-Regular.ttf", "PlayfairDisplay-Bold.ttf"),
    'RobotoSlab': _FontFiles("RobotoSlab-Regular.ttf", "RobotoSlab-Bold.ttf"),
  };

  _FontFiles get _fallbackFont => this._fontFilesMap['OpenSans']!;

  _FontFiles _getFontOrFallback(String fontName) {
    _FontFiles? fontFiles = _fontFilesMap[fontName];
    if (fontFiles == null) {
      print(
          'font ${fontName} was not found, using the default font as fallback');
      return _fallbackFont;
    }
    return fontFiles;
  }
}
