from selenium import webdriver
driver = webdriver.Firefox()

driver.get("https://play.google.com/store/apps/details?id=info.blockchain.merchant")

driver.find_element_by_css_selector("button.content.id-view-permissions-details.fake-link").click()

# for whole page
htmlElem = driver.find_element_by_tag_name('html')
print htmlElem.text

# OR for just the permissions text

permdeets = driver.find_element_by_css_selector("div.permission-buckets.id-permission-buckets")
print permdeets.text
