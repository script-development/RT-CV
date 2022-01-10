import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'utils.dart';

class ListWithHeader {
  ListWithHeader(this.icon, this.title, this.widgets);

  final IconData icon;
  final String title;
  final List<Widget> widgets;

  get length => widgets.length;
}

class WrapLayoutBlock extends StatelessWidget {
  WrapLayoutBlock(this.listWithHeader, this.headerColor);

  final ListWithHeader listWithHeader;
  final BgColor headerColor;

  @override
  Widget build(Context context) {
    List<List<Widget>> columns = [];
    for (var i = 0; i < listWithHeader.widgets.length; i++) {
      int idxInRow = i % 3;
      Widget toAdd = Expanded(
        child: Padding(
          padding: EdgeInsets.only(
            left: idxInRow == 0 ? 0 : 2,
            right: idxInRow == 2 ? 0 : 2,
            top: i > 2
                ? 4
                : 0, // Only add padding to top if we are on the second row or more
          ),
          child: listWithHeader.widgets[i],
        ),
      );
      if (idxInRow == 0) {
        columns.add([
          Expanded(child: toAdd),
          Expanded(child: Container()),
          Expanded(child: Container()),
        ]);
      } else {
        columns[columns.length - 1][idxInRow] = toAdd;
      }
    }

    return LayoutBlockBase(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ListTitle(
            listWithHeader.icon,
            listWithHeader.title,
            headerColor,
          ),
          Column(
            children: columns
                .map((row) => Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: row,
                    ))
                .toList(),
          ),
        ],
      ),
    );
  }
}

class ColumnsLayoutBlock extends StatelessWidget {
  ColumnsLayoutBlock(this.widgets, this.headerColor);

  final List<ListWithHeader> widgets;
  final BgColor headerColor;

  @override
  Widget build(Context context) {
    List<Widget> left = [];
    List<Widget> right = [];

    for (int i = 0; i < widgets.length; i++) {
      ListWithHeader widget = widgets[i];
      ListTitle title = ListTitle(
        widget.icon,
        widget.title,
        headerColor,
        margin: i > 1 ? const EdgeInsets.only(top: 10) : null,
      );

      if (i % 2 == 0) {
        left.add(title);
        left.addAll(widget.widgets);
      } else {
        right.add(title);
        right.addAll(widget.widgets);
      }
    }

    return LayoutBlockBase(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(right: PdfPageFormat.cm / 2),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: left,
              ),
            ),
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(left: PdfPageFormat.cm / 2),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: right,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class LayoutBlockBase extends StatelessWidget {
  LayoutBlockBase({required this.child});

  final Widget child;

  @override
  Widget build(Context contex) {
    return Padding(
      padding: const EdgeInsets.only(
        left: PdfPageFormat.cm,
        right: PdfPageFormat.cm,
        top: PdfPageFormat.cm / 2,
      ),
      child: child,
    );
  }
}

class ListTitle extends StatelessWidget {
  ListTitle(this.icon, this.title, this.headerColor, {this.margin});

  final IconData icon;
  final String title;
  final BgColor headerColor;
  final EdgeInsets? margin;

  @override
  Widget build(Context context) {
    return Container(
      margin: margin,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            color: headerColor.bgColor,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.only(right: 5),
                    child: Icon(
                      icon,
                      size: 10,
                      color: headerColor.textColor,
                    ),
                  ),
                  Text(
                    title,
                    overflow: TextOverflow.clip,
                    style: TextStyle(
                      color: headerColor.textColor,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),
          ),
          Container(
            width: double.infinity,
            height: 2,
            color: headerColor.bgColor,
          ),
        ],
      ),
    );
  }
}
