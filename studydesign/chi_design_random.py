import random
import csv

neworder = []

# create 25 questions in random order from group a, without repeating interface or category
def makesequence1():
	oldorder = []
	input_file = csv.DictReader(open("group1.csv"))
	for row in input_file:
		oldorder.append(row)
	seed = random.choice(oldorder)
	neworder.append(seed)
	oldorder.remove(seed)
	i = 0
	while i < 24:
		choice = random.choice(oldorder)
		previous = neworder[-1]
		if (previous['interface'] == choice['interface']) or (previous['category'] == choice['category']):
			print 'found repeated order :( try again\n'
		else:
			neworder.append(choice)
			print "added to neworder"
			oldorder.remove(choice)
			print "removed from oldorder"
			i+=1
	for item in neworder:
		item['position'] = neworder.index(item)
	return neworder

group1order = makesequence1()

neworder = []

newcsv = open('newa.csv', 'w')
for item in group1order:
	newcsv.write("%s,%s,%s,%s,%s\n" % (item['category'], item['appair'], item['interface'], item['group'], item['position']))

# create 25 questions in random order from group b, without repeating interface or category

def makesequence2():
	oldorder = []
	input_file = csv.DictReader(open("group2.csv"))
	for row in input_file:
		oldorder.append(row)
	seed = random.choice(oldorder)
	neworder.append(seed)
	oldorder.remove(seed)
	i = 0
	while i < 24:
		choice = random.choice(oldorder)
		previous = neworder[-1]
		for item in group1order:
			if (item['position'] == choice['position']) and ((item['appair'] == choice['apppair']) and (item['interface'] == choice['interface'])):
				print 'found repeated position'
				choice = random.choice(oldorder)
		if (previous['interface'] == choice['interface']) or (previous['category'] == choice['category']):
			print 'found repeated order :( try again\n'
		else:
			neworder.append(choice)
			print "added to neworder"
			oldorder.remove(choice)
			print "removed from oldorder"
			i+=1
	for item in neworder:
		item['position'] = neworder.index(item)
	return neworder

group2order = makesequence2()

neworder = []

newcsv = open('newb.csv', 'w')
for item in group2order:
	newcsv.write("%s,%s,%s,%s,%s\n" % (item['category'], item['appair'], item['interface'], item['group'], item['position']))

# create 25 questions in random order from group c, without repeating interface or category

def makesequence3():
	oldorder = []
	input_file = csv.DictReader(open("group3.csv"))
	for row in input_file:
		oldorder.append(row)
	seed = random.choice(oldorder)
	neworder.append(seed)
	oldorder.remove(seed)
	i = 0
	while i < 24:
		choice = random.choice(oldorder)
		previous = neworder[-1]
		for item in group2order:
			if (item['position'] == choice['position']) and ((item['appair'] == choice['apppair']) and (item['interface'] == choice['interface'])):
				print 'found repeated position'
				choice = random.choice(oldorder)
		if (previous['interface'] == choice['interface']) or (previous['category'] == choice['category']):
			print 'found repeated order :( try again\n'
		else:
			neworder.append(choice)
			print "added to neworder"
			oldorder.remove(choice)
			print "removed from oldorder"
			i+=1
	for item in neworder:
		item['position'] = neworder.index(item)
	return neworder

group3order = makesequence3()

neworder = []

newcsv = open('newc.csv', 'w')
for item in group3order:
	newcsv.write("%s,%s,%s,%s,%s\n" % (item['category'], item['appair'], item['interface'], item['group'], item['position']))

def makesequence4():
	oldorder = []
	input_file = csv.DictReader(open("group4.csv"))
	for row in input_file:
		oldorder.append(row)
	seed = random.choice(oldorder)
	neworder.append(seed)
	oldorder.remove(seed)
	i = 0
	while i < 24:
		choice = random.choice(oldorder)
		previous = neworder[-1]
		for item in group2order:
			if (item['position'] == choice['position']) and ((item['appair'] == choice['apppair']) and (item['interface'] == choice['interface'])):
				print 'found repeated position'
				choice = random.choice(oldorder)
		if (previous['interface'] == choice['interface']) or (previous['category'] == choice['category']):
			print 'found repeated order :( try again\n'
		else:
			neworder.append(choice)
			print "added to neworder"
			oldorder.remove(choice)
			print "removed from oldorder"
			i+=1
	for item in neworder:
		item['position'] = neworder.index(item)
	return neworder

group4order = makesequence4()

neworder = []

newcsv = open('newd.csv', 'w')
for item in group4order:
	newcsv.write("%s,%s,%s,%s,%s\n" % (item['category'], item['appair'], item['interface'], item['group'], item['position']))

def makesequence5():
	oldorder = []
	input_file = csv.DictReader(open("group5.csv"))
	for row in input_file:
		oldorder.append(row)
	seed = random.choice(oldorder)
	neworder.append(seed)
	oldorder.remove(seed)
	i = 0
	while i < 24:
		choice = random.choice(oldorder)
		previous = neworder[-1]
		for item in group2order:
			if (item['position'] == choice['position']) and ((item['appair'] == choice['apppair']) and (item['interface'] == choice['interface'])):
				print 'found repeated position'
				choice = random.choice(oldorder)
		if (previous['interface'] == choice['interface']) or (previous['category'] == choice['category']):
			print 'found repeated order :( try again\n'
		else:
			neworder.append(choice)
			print "added to neworder"
			oldorder.remove(choice)
			print "removed from oldorder"
			i+=1
	for item in neworder:
		item['position'] = neworder.index(item)
	return neworder

group5order = makesequence5()

neworder = []

newcsv = open('newe.csv', 'w')
for item in group5order:
	newcsv.write("%s,%s,%s,%s,%s\n" % (item['category'], item['appair'], item['interface'], item['group'], item['position']))

