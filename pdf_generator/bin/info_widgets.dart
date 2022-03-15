import 'package:pdf/pdf.dart';
import 'package:pdf/widgets.dart';

import 'cv.dart';
import 'utils.dart';

class ClientInfo extends StatelessWidget {
  ClientInfo({
    required this.personalInformation,
    this.driversLicenses,
  });

  List<Widget> children = [];
  final PersonalDetails personalInformation;
  final List<String>? driversLicenses;

  final TextStyle labelStyle = TextStyle(
    fontSize: 10,
    color: PdfColors.grey800,
  );
  final TextStyle valueStyle = TextStyle(
    fontSize: 10,
    color: PdfColors.black,
  );

  tryAddToList(String label, String? value) {
    if (value != null) {
      children.add(
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(label + ": ", style: labelStyle),
            Text(
              value,
              overflow: TextOverflow.clip,
              style: valueStyle,
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(Context context) {
    children = [];
    tryAddToList("E-mail", personalInformation.email);
    tryAddToList("Telefoon", personalInformation.phoneNumber);
    if (driversLicenses != null) {
      switch (driversLicenses!.length) {
        case 0:
          // Do not add the drivers licenses
          break;
        case 1:
          tryAddToList("Rijbewijs", driversLicenses![0]);
          break;
        default:
          tryAddToList("Rijbewijzen", driversLicenses!.join(", "));
      }
    }

    final EdgeInsets verticalPadding = EdgeInsets.symmetric(vertical: 20);

    if (!personalInformation.hasAddress) {
      if (personalInformation.zip != null) {
        String? postalCodePlace =
            guessPostalCodePlace(personalInformation.zip!);

        if (postalCodePlace != null)
          tryAddToList("Postcode",
              "${personalInformation.zip} (regio ${postalCodePlace})");
        else
          tryAddToList("Postcode", personalInformation.zip);
      }

      return Padding(
        padding: verticalPadding,
        child: Wrap(children: children, spacing: 10),
      );
    }

    return Padding(
      padding: verticalPadding,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            constraints: BoxConstraints(
              minWidth: 150,
            ),
            child: Padding(
              padding: EdgeInsets.only(right: 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text("Plaats: ", style: labelStyle),
                      Text(personalInformation.city!, style: valueStyle),
                    ],
                  ),
                  Row(
                    children: [
                      Text("Adres: ", style: labelStyle),
                      Text(
                          "${personalInformation.streetName} ${personalInformation.houseNumber} ${personalInformation.houseNumberSuffix}",
                          style: valueStyle),
                    ],
                  ),
                  Row(
                    children: [
                      Text("Postcode: ", style: labelStyle),
                      Text(personalInformation.zip!, style: valueStyle),
                    ],
                  ),
                ],
              ),
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: children,
          ),
        ],
      ),
    );
  }
}

class WorkExpWidget extends StatelessWidget {
  WorkExpWidget(WorkExperience this.exp);

  final WorkExperience exp;

  @override
  Widget build(Context context) {
    return Padding(
      padding: EdgeInsets.only(bottom: 10),
      child: ListEntry(
        exp.profession ?? '??',
        company: exp.employer,
        description: exp.description,
        from: exp.startDate,
        to: exp.endDate,
      ),
    );
  }
}

class EducationWidget extends StatelessWidget {
  EducationWidget(Education this.education);

  final Education education;

  @override
  Widget build(Context context) {
    return Padding(
      padding: EdgeInsets.only(bottom: 10),
      child: ListEntry(
        education.name,
        company: education.institute,
        from: education.startDate,
        to: education.endDate,
        description: education.description,
      ),
    );
  }
}

class ListEntry extends StatelessWidget {
  ListEntry(
    this.title, {
    this.description,
    this.company,
    this.from,
    this.to,
  });

  final String title;
  final String? company;
  final String? description;
  final DateTime? from;
  final DateTime? to;

  @override
  Widget build(Context context) {
    List<Widget> children = [
      Flexible(
        child: Text(
          title,
          overflow: TextOverflow.clip,
          style: TextStyle(
            fontSize: 10,
            fontWeight: FontWeight.bold,
            lineSpacing: 2,
          ),
        ),
      ),
    ];

    TextStyle contentStyle = TextStyle(
      fontSize: 10,
      color: PdfColors.grey800,
      lineSpacing: 2,
    );
    TextStyle dateStyle = TextStyle(
      fontSize: 8,
      color: PdfColors.grey600,
      lineSpacing: 2,
    );

    if (company != null && company!.isNotEmpty) {
      children.add(
        Row(children: [
          Text(
            company!,
            overflow: TextOverflow.clip,
            style: contentStyle,
          ),
        ]),
      );
    }

    String? fromStr = formatDate(from);
    String? toStr = formatDate(to);
    if (fromStr != null || toStr != null) {
      if (toStr == null || fromStr == null) {
        children.add(
          Row(children: [
            Text(fromStr ?? toStr ?? '??', style: dateStyle),
          ]),
        );
      } else {
        children.add(
          Row(children: [
            Text(fromStr, style: dateStyle),
            Text(" - ", style: dateStyle),
            Text(toStr, style: dateStyle),
          ]),
        );
      }
    }

    if (description != null && description!.isNotEmpty) {
      children.add(Flexible(
        child: Text(
          (description!.length > 300)
              ? description!.substring(0, 300) + '..'
              : description!,
          overflow: TextOverflow.clip,
          style: contentStyle,
        ),
      ));
    }

    return Padding(
      padding: const EdgeInsets.only(top: 5),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: children,
      ),
    );
  }
}
