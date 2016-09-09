

// textutil -convert txt *

var parse = require('csv-parse/lib/sync'),
	csvstr = require('csv-stringify'),
	fs = require('fs'),
	promise = require('bluebird'),
	_ = require('lodash'),
	headers,
	qs = require('querystring'),
	config = JSON.parse(fs.readFileSync('./munge-config.json')),
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
				d[fname.slice(0,-'.txt'.length)] = load_transcript([srcdir,fname].join('/'));
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
					rounds = filed;

				if (pid !== fname.slice(0,-'.json'.length)) { 
					console.warn('file doestn match participant id ', pid, fname);
				}
				d[pid] = rounds;
				return d;
			}, 
			{});
}, gen_out = (transcripts, rounds, fakeapps) => {
	// get fields x
	var field_names = [
		'id',
		'round',
		'participant',
		'condition',
		'domain',
		'app_a',
		'type_a',
		'app_b',			
		'type_b',
		'chosen',
		'type_chosen',
		'elapsed_secs',			
		'confidence',
		'thinkaloud'
	],
	field_values = [
		(rounds, r, ri) => [rounds.participant,''+ri].join('-'), // unique id
		(rounds, r, ri) => ri+1,
		(rounds, r, ri) => rounds.participant,
		(rounds, r, ri) => r.cond,
		(rounds, r, ri) => r.domain,
		(rounds, r, ri) => r.a,
		(rounds, r, ri) => fakeapps[r.a],
		(rounds, r, ri) => r.b,
		(rounds, r, ri) => fakeapps[r.b],
		(rounds, r, ri) => r.result.chosen,
		(rounds, r, ri) => fakeapps[r.result.chosen],
		(rounds, r, ri) => Math.round(r.result.elapsed/1000.0),
		(rounds, r, ri) => parseInt(r.result.confidence.slice('likert'.length+1)),
		(rounds, r, ri) => transcripts[rounds.participant] && transcripts[rounds.participant][ri] || '~'
	];

	var rows = [field_names].concat(_.flatten(_.keys(rounds).map((participant) => {
		var rdata = rounds[participant];
		return _(rdata.rounds).map((r,i) => field_values.map((f) => {
			console.info(rdata.participant, i, field_names[field_values.indexOf(f)], f(rdata,r,i));
			return f(rdata,r,i);
		})).value();
	})));

	return new Promise((acc, rej) => {
		csvstr(rows, (err, output) => {
			if (!err) { return acc(output); }
			console.error("Error ", err);
			rej(err);
		});
	});
}, loadCSV = (fname) => {
	var text = 	fs.readFileSync(fname).toString();
	console.log("Parsing file ", fname, "(", text.length, ")");
	var data = parse(text, {max_limit_on_data_read:9999999999});
	headers = data[0];
	data = data.slice(1);
	data = data.map((x) => _.zipObject(headers,x));
	return data;
}, app2type = (apparr) => {
	// simply takes array [ { "App type code":'', "fake name": .. } ]
	// -> { fake name : app type code }

	return apparr.reduce((d, app) => {
		d[app['fake name']] = app['App type code'];
		return d;
	}, {});
},
main = (mode) => {
	var ts = load_transcripts(), 
		rs = load_rounds(),
		fakeapps = app2type(loadCSV(config.fakeapps)),
		fout = config.out;

	if (!fout) { 
		console.error("No output directory specified, please set out_dir in qual-chop-config");
		return;
	}

	console.info('loaded ', _.keys(ts).length, ' transcripts', _.keys(ts));
	console.info('loaded ', _.keys(rs).length, ' rounds ', _.keys(rs));

	console.info('fake apps ', fakeapps);


	gen_out(ts, rs, fakeapps).then((output) => {
		console.info("Writing output ", fout, fout.length);
		fs.writeFileSync(fout, output);
		console.log('done.');
	}).catch((err) => {
		console.error("Error, terminating ", err);
	});
};

if (require.main === module) { main(); }