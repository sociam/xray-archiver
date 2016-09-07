var parse = require('csv-parse/lib/sync'),
	fs = require('fs'),
	_ = require('lodash'),
	headers,
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./qual-chop-config.json')),
	round_re = /ROUND\W+(\d+)/;

var load_transcript = (fname) => {
	var text = fs.readFileSync(fname).toString(),
		first_sep_match = text.match(round_re),
		fidx = first_sep_match && text.indexOf(first_sep_match[0]),
		rtext = fidx && text.slice(fidx).split('\n').filter((x) => x.length > 0),
		rounds = [],
		cur_r = first_sep_match && parseInt(first_sep_match[1]);

	if (!first_sep_match) { 
		console.error("No separator found in ", fname, " skipping");
		return;
	}

	rtext.map((line) => {
		var rT = rounds[cur_r] || '',
			m = line.match(round_re);
		if (m) {
			// incr and skip
			cur_r = parseInt(m[1]);
			console.info('next round ', fname, ' r ', cur_r);			
			return;
		}
		rounds[cur_r] = rT + line;
	});

	return rounds;
}, 
load_transcripts = () => {
	var srcdir = config.transcripts;
	return fs.readdirSync(srcdir)
		.filter((fname) => fname.indexOf('.txt') >= 0)
		.reduce((d,fname) => 
			{ 
				d[fname] = load_transcript([srcdir,fname].join('/'));
				return d;
			}, {});
},
load_rounds = () => {
	// loads data files in batch
	var srcdir = config.rounds;
	return fs.readdirSync(srcdir)
		.filter((fname) => fname.indexOf('.json') >= 0)
		.reduce((d,fname) => 
			{ 
				var filen = [srcdir,fname].join('/'),
					filed = JSON.parse(fs.readFileSync(filen).toString()),
					pid = filed.participant,
					rounds = filed.rounds;

				if (pid !== fname.slice(0,-'.json'.length)) { 
					console.warn('file doestn match participant id ', pid, fname);
				}
				d[pid] = rounds;
				return d;
			}, 
			{});
}, by_field = (field, transcripts, rounds) => {
	// get fields x
	var out_dir = config.out_dir,
		vals = _(rounds).values().flatten().map((x) => x[field]).uniq().value();

	if (!out_dir) { 
		console.error("No output directory specified, please set out_dir in qual-chop-config");
		return;
	}

	console.log('vals :', vals);

}, main = (mode) => {
	var ts = load_transcripts(), rs = load_rounds();
	if (mode === undefined || mode == 'cond') {
		by_field('cond', ts, rs);
	}
};
// pdci []
if (require.main === module) { 
	// console.info(process.argv.length);
		//   console.log(index + ': ' + val);
	// });	
	if (process.argv.length === 2 || process.argv.length === 3) { 
		// console.log('fo ', process.argv.length);
		return main(process.argv[2]);
	}  else {
		main(); 
	}
}