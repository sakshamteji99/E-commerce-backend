# E-commerce-backend

README.MD

This is a quick illustration to a dummy e-commerce website namely Monk Commerce. The project is a totally backend development effort with no UI, no database and only certain API endpoints which mainly revolve around the functionality on managing discounts. 

#####How to Run: 

1) Download go
2) Save the code
3) Run command in the directory: go run .

#####Assumptions and Limitations: 

1) As no database is used, everything is stored in the code, hence once the server is turned off you will loose all your data. This is done to reduce the development effort
2) I Introduced a points concept instead of giving the user straight cash discount. These points could be converted to cash with some logic which would then be transferred to Monk Commerce Wallet. To make this more appealing, we can give the customer an option between cash discount and 1.5*(cash discount) equivalent of points. This will tempt the customer to get more points(feels more lucrative) which again customer will use to buy more products ultimaely making it a profitable deal for both user and Monk Commerce. The customer will accumulate points for every purchase.
3) Points are calculated as 1.2x the total cash discount.
4) Logic is implemented more in the terms of percentage of cart value and not flat off. This is to minimise the hardcoding of discount values. Percentage discouts appear more realistic and work wonders with any cart value.
5) Free delivery applies only for the first 5 orders and if the cart value exceeds Rs 1000
6) Special sale days are fixed as: Christmas (December 25) && Independence Day (August 15). Also company anniversary discounts are fixed on December 1.
7) Concurrent updates to shared maps like products and customers may lead to race conditions which is not handled to avoid further code complexity
8) Admin secret (Admin123) is predefined. No proper authentication and authorization implementation is done. 
9) No logging monitoring and alerting set up is done. 
10) No testing is done at all for the service
11) Coupons cannot expire 


API Endpoints:

#####Admin Endpoints:
1) Set Admin: POST /admin/set-admin
Body:
{
  "secret": "Admin123"
}

2) Add or Remove Product: POST /admin/add-or-remove/product
Body:
{
  "secret": "Admin123",
  "action": "add",
  "name": "Product Name",
  "quantity": 10,
  "sale_price": 100.0
}

3) For removing a product: POST /admin/add-or-remove/product
Body:
{
  "secret": "Admin123",
  "action": "remove",
  "product_id": 1
}

4) Add or Remove Coupon Category: POST /admin/add-or-remove/coupon/category
Body:
{
  "action": "add",
  "name": "Category Name"
}

5) Add or Remove Coupon: POST /admin/add-or-remove/coupon/
Body:
{
  "secret": "Admin123",
  "action": "add",
  "category": "Category Name",
  "min_value": 500,
  "max_value": 2000,
  "discount_percent": 10,
  "max_discount": 100,
  "description": "My Description"
}

#####Customer Endpoints

1) Check Discounts: POST /customer/check-discounts
Body:
{
  "cart_value": 1000,
  "customer_id": 1
}

2) Apply Discount: POST /customer/apply-discount
Body:
{
  "cart_value": 1000,
  "coupon_id": 1,
  "customer_id": 1,
  "option": "cash"
}

3) Place Order: POST /customer/place-order
Body:
{
  "customer_id": 1,
  "product_id": 1,
  "quantity": 2,
  "cart_value": 1000
}

4) Add Money to Wallet: POST /customer/add-money-to-wallet
Body:
{
  "customer_id": 1,
  "amount": 500,
  "use_points": false
}

##### Features and Discount Cases implemented: 

1) Cart Value-Based Coupons:
Coupons are applied if the cart value falls within a defined range (min_value and max_value), with a specific discount percentage and a maximum discount cap.

2) Special Sale Day Discount:
An additional fixed percentage discount (extraDiscountPercent) is automatically applied if the sale day logic is triggered.

3) Points would be accumulated on every purchase and later could be converted to wallet money. 

4) Birthday Discount and Anniversary Discount: If the customer’s birthday month matches the current month, an additional discount is added and 
a special flat discount for all customers on December 1.




#####Potential Discount Logics which were thought of but weren't implemented: 

1) Product Specific Discounts(BxGy)
Complexity:
Needed complex mappings and validations and could conflict with other cart-level discounts. Also, logics had to be intensely hardcoded for this and for every cart value, it was a complex logic. 

2) Time Sensitive Discounts: E.g., After item is added in cart, offer 10% more discount if bought within X hours or discounts of 10% off between 3 PM and 6 PM.
Complexity:
Requires handling of precise server time and timezone differences for different users.

3) Referral Based Discounts:
Complexity:
Requires tracking and managing of more complex logic for type "User".

4) Loyalty Program Discounts: The more worth of items purchased, the more points added. 

5) Subscription-Based Discounts

6) Location-Based Discounts

For any further queries, reach out to sakshamteji99@gmail.com or call at 9501151175