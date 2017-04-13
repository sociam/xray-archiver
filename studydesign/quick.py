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

order = makesequence1()

group = 'a.'

sequence = ''

for item in order:
	sequence+= "%s::%s::%s," % (item['appair'], item['category'], item['interface'])

print group+sequence
"""
newcsv = open('new.csv', 'w')
newcsv.write(group +'.')
for item in grouporder:
	newcsv.write("%s::%s::%s," % (item['appair'], item['category'], item['interface']))




"""