# Herschel
The Google Spreadsheet Writing / Reading Library for Golang.

## Features
* Read table from Google Spreadsheet
* Write table to Google Spreadsheet
* Cell decoration and formatting
* Configure sheet

## Installation
```
go get github.com/yokoe/herschel
```

## Quick Start
You need service account credentials or user account credentials for api call. 
https://developers.google.com/sheets/api/guides/authorizing


```
spreadsheetID := "1234567890...SpreadsheetID"

client, err := herschel.NewClient(option.WithServiceAccountCredentials("service-account.json"))
if err != nil {
    log.Fatalf("invalid config or token %s", err)
}

// Read data from Spreadsheet
table, err := client.ReadTable(spreadsheetID, "Sheet 1")
if err != nil {
    log.Fatalf("read error %s", err)
}
fmt.Print(table.GetValue(0, 0))

// Edit data
table.PutValue(0, 0, "Hello from Herschel")

// Write to Spreadsheet
err = client.WriteTable(spreadsheetID, "Sheet 1", table)
if err != nil {
    log.Fatalf("write error %s", err)
}
```

## How to use
### Writing table
```
client, err := ...

table := NewTable(2, 2)
table.Put(0, 0, "Updated: "+time.Now().Format("2006-01-02 03:04"))
table.Put(1, 0, 1234567890)
table.Put(0, 1, "fuga")
table.Put(1, 1, 0.2530)

table.SetBackgroundColor(0, 0, color.Black)
table.SetBackgroundColor(1, 1, color.RGBA{128, 0, 0, 0})
table.SetNumberFormatPattern(1, 0, "#,###")
table.SetNumberFormatPattern(1, 1, "#.00%")

table.FrozenRowCount = 1
table.FrozenColumnCount = 2

err = client.WriteTable(spreadsheetID, "Sheet 1", table)
if err != nil {
    // Error handling
}
```

### Reading table
```
client, err := ...

table, err := client.ReadTable(spreadsheetID, "Sheet 1")
if err != nil {
    // Error handling
}

// table.GetValue(0, 0)
```

### Table manipulation
#### Get / Put
```
table.PutValue(0, 0, "Hello world")
table.PutValuesAtRow(1, "Hello", "World")

table.GetValue(0, 0) // "Hello world"
table.GetStringValue(0, 0) // "Hello world"
table.GetValuesAtRow(1) // "Hello", "World"
```

#### Finding row
```
table.PutValuesAtRow(0, "a", "b", "c")
table.PutValuesAtRow(1, "d", "e", "f")
table.PutValuesAtRow(2, "g", "h", "i")

table.IndexOfRowWithPrefix("d", "e") // 1
table.IndexOfRowWithPrefix("d") // 1
table.IndexOfRowWithPrefix("b") // -1

```

#### Freeze rows / cols
```
table.FrozenRowCount = 1
table.FrozenColumnCount = 2
```

#### Removing row / col
```
table.RemoveRowAtIndex(3)
```

### Sheet manipulation
```
client.AddSheet("spreadsheetID", "NewSheet")
client.DeleteSheet("spreadsheetID", "NewSheet")
client.RecreateSheet("spreadsheetID", "NewSheet")
```

### Worksheet manipulation
```
id, err := client.CreateNewSpreadsheet(config, token, "NewWorksheet")
```

## Authentication
You can authenticate using service accounts or user accounts.

```
client, err := NewClient(option.WithServiceAccountCredentials(credentialsFilePath))
client, err := NewClient(option.WithConfigFileAndTokenFile(configFilePath, userCredentialsFilePath))
```

## Development
### Run testcases with api call
Testcases require api access to spreadsheet will be skipped in default.
To run all testcases, please set service account credentials json file to `SPREADSHEET_CREDENTIAL_FILE`.

```
SPREADSHEET_CREDENTIAL_FILE=/path/to/credentials.json go test . -v -cover
```
