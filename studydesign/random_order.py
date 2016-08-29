import random
import csv

oldorder = []
neworder = []

input_file = csv.DictReader(open("group1.csv"))
for row in input_file:
	oldorder.append(row)

seed = random.choice(oldorder)

neworder.append(seed)

print neworder

i = 0
while i < 25:
	choice = random.choice(oldorder)
	previous = neworder[-1]
	if ((previous['interface'] == choice['interface']) or (previous['category'] == choice['category'])):
		print 'found repeated order!!!!!!!!!!!!\n'
	else:
		neworder.append(choice)
		i+=1

print neworder

print "checking for repeated interfaces"
for choice1, choice2 in zip(neworder, neworder[1:]):

	if choice1['interface'] == choice2['interface']:
		print 'found one!\n'

print "checking for repeated categories"
for choice1, choice2 in zip(neworder, neworder[1:]):

	if choice1['category'] == choice2['category']:
		print 'found one!\n'