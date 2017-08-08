# Suggester - Alternative App Suggestions

This script scrapes [AlternativeTo](www.alternativeto.com) for alternative app suggestions.

The implementation will be as follows:

```
For Each (AppID, Title) in the Database
    fetch HTML from www.alternativeto.com/browse/search?q=<title>
    if the title matched directly
        For Each suggested alternative apps
            Log them to the Database
    Else do something else to get suggestions...
```
