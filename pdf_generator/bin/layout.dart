import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'utils.dart';
import 'args.dart';

class ListWithHeader {
  ListWithHeader(this.icon, this.title, this.widgets);

  final IconData icon;
  final String title;
  final List<Widget> widgets;

  get length => widgets.length;
}

class WrapLayoutBlock extends StatelessWidget {
  WrapLayoutBlock(this.listWithHeader, this.style);

  final ListWithHeader listWithHeader;
  final Style style;

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
      child: Container(
        decoration: BoxDecoration(
          border: style.layoutStyle == LayoutStyle.style_3
              ? Border(
                  left: BorderSide(
                    color: style.subHeaderBackgroundColor,
                    width: 2,
                  ),
                )
              : null,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ListTitle(
              listWithHeader.icon,
              listWithHeader.title,
              style,
            ),
            Padding(
              padding: EdgeInsets.only(
                  left: style.layoutStyle == LayoutStyle.style_3 ? 10 : 0),
              child: Column(
                children: columns
                    .map((row) => Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: row,
                        ))
                    .toList(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ColumnsLayoutBlock extends StatelessWidget {
  ColumnsLayoutBlock(this.widgets, this.style);

  final List<ListWithHeader> widgets;
  final Style style;

  @override
  Widget build(Context context) {
    List<Widget> left = [];
    List<Widget> right = [];

    for (int i = 0; i < widgets.length; i++) {
      ListWithHeader widget = widgets[i];

      Widget listEntriesWidget = Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: widget.widgets,
      );

      Widget listWithHeaderToadd = Padding(
        padding: EdgeInsets.only(top: i > 1 ? 10 : 0),
        child: Container(
          decoration: BoxDecoration(
            border: style.layoutStyle == LayoutStyle.style_3
                ? Border(
                    left: BorderSide(
                      color: style.subHeaderBackgroundColor,
                      width: 2,
                    ),
                  )
                : null,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ListTitle(
                widget.icon,
                widget.title,
                style,
              ),
              Padding(
                padding: EdgeInsets.only(
                    left: style.layoutStyle == LayoutStyle.style_3 ? 10 : 0),
                child: listEntriesWidget,
              ),
            ],
          ),
        ),
      );

      if (i % 2 == 0) {
        left.add(listWithHeaderToadd);
      } else {
        right.add(listWithHeaderToadd);
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
  ListTitle(this.icon, this.title, this.style);

  final IconData icon;
  final String title;
  final Style style;

  @override
  Widget build(Context context) {
    var children = [
      Container(
        decoration: BoxDecoration(
          borderRadius: style.layoutStyle == LayoutStyle.style_2
              ? BorderRadius.circular(5)
              : null,
          color: style.subHeaderBackgroundColor,
        ),
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
                  color: style.subHeaderTextColor,
                ),
              ),
              Text(
                title,
                overflow: TextOverflow.clip,
                style: TextStyle(
                  color: style.subHeaderTextColor,
                  fontSize: 10,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
      ),
    ];

    if (style.layoutStyle == LayoutStyle.style_1) {
      children.add(Container(
        width: double.infinity,
        height: 2,
        color: style.subHeaderBackgroundColor,
      ));
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: children,
    );
  }
}
