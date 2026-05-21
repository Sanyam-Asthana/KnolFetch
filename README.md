# KnolFetch

Knolfetch is a CLI tool which fetches the latest research in a specified field (using arXiv API) and dumps it in the terminal. I originally made this for myself, as a more useful alternative to fastfetch (no hate).

## Building

1. Clone the repository
2. Run `go build ./build/knolfetch`
3. The compiled binary will be in `/build`

## Usage

1. Create the file `~/.config/knolfetch/config.xml` and directory `~/.cache/knolfetch/`.
2. Configure the settings in the `config.xml` file.
3. Run the pre-compiled binary, or run from source using `go run`, with the flag `--fetch`. This will make the request to arXiv.org and save the results in the cache.
4. Run the program without any flags to see the results.

## Configuration

### Example `config.xml`

```
<config>
        <category>
                <name>quant-ph</name>
                <maxresults>2</maxresults>
                <maxabstractlength>100</maxabstractlength>
        </category>

        <category hideauthor="1">
                <name>cs.CL</name>
                <maxresults>2</maxresults>
                <maxabstractlength>100</maxabstractlength>
        </category>
</config>
```

### Options

- `<category>` block contains the configuration for a specific category.

- `<name>` blocks specify the field of the research, according to the subject codes of arXiv.org.

- Setting the attributes `hideauthor`, `hidetitle`, `hideabstract`, `hidelink` of the `<category>` tag hide the respective elements.

- `<maxresults>` specifies the maximum number of top results to be returned of the specific category.

- `<maxabstractlength>` specifies the maximum length of the abstract.

### Regular fetching

A cronjob can be set for fetching the latest research using `knolfetch --fetch` at a set interval easily.

**Example crontab:**

```
*/60 * * * * knolfetch --fetch
```

The above crontab configuration runs the command every 60 minutes.
