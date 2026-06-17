## سناریو ۱: تعریف قانون کمیسیون

**نقش:** مدیر سیستم

**شرح:**
سیستم برای هر فروش از فروشنده کمیسیون دریافت می‌کند. مدیر سیستم باید قوانین کمیسیون را برای هر کالا یا دسته‌بندی تعریف کند.

## senario 1̱: tarif ghanon kamision **naqsh:** madir sistam **sharh:** sistam baraye npar forosh az foroshandeh kamision daryaft mikand. madir sistam bayad ghavanin kamision ra baraye npar kala ya dastehbandi tarif kand. npar kala faghat yek ghanon kamision mitavand dashteh bashod. **jarian esli:** 1. madir kalaye moord nazar ra entekhab mikand 2. darsad kamision ve sharayet an ra taein mikand 3. sistam ghanon kamision ra dar jodol `commission` sabat mikand **masal:** `kalaye A: - madel forosh: retail - darsad kamision: 1̱0̱% - bazeh ghimet: 5̱0̱,0̱0̱0̱ ta 5̱0̱0̱,0̱0̱0̱ toman - hadaghal taedad: 1̱` ## senario 2̱: mohasbeh kamision forosh **naqsh:** sistam **sharh:** bad az npar forosh mofegh, sistam bayad sonpam khod ra az forosh mohasbeh kand. moghdar kamision bar asas darsad tarif shodeh ve sharayet ghanon mohasbeh mishod. **jarian esli:** 1. yek forosh enjam mishod 2. sistam ghanon kamision marbut bah kala ra pida mikand 3. sistam sharayet ghanon ra barresi mikand: - aya ghimet forosh dar bazeh `min_price` ta `max_price` gharar dard? - aya taedad forosh rafteh hadaghal `min_qty` ra dard? 4. eger sharayet bargharar bashod, kamision mohasbeh mishod: `kamision = mobolgh forosh * (rate_percent / 100)`

بیشتر نشان دادن
۱٬۰۶۹

## Scenario 1: Defining a Commission Rule

**Role:** System Administrator

**Description:**
The system receives a commission from the seller for each sale. The system administrator must define commission rules for each product or category. Each product can have only one commission rule.

**Main flow:**

1. The administrator selects the desired product
2. Determines the commission percentage and its conditions
3. The system records the commission rule in the `commission` table

**Example:**

```
Product A:
- Sales model: retail
- Commission percentage: 10%
- Price range: 50,000 to 500,000 Tomans
- Minimum quantity: 1
```

## Scenario 2: Calculating sales commission

**Role:** System

**Description:**
After each successful sale, the system must calculate its share of the sale. The commission amount is calculated based on the defined percentage and the conditions of the rule.

**Main flow:**

1. A sale is made
2. The system finds the commission rule for the item
3. The system checks the rule conditions:

- Is the sale price between `min_price` and `max_price`?
- Is the quantity sold at least `min_qty`?

4. If the conditions are met, the commission is calculated:

```
Commission = Sale Amount * (rate_percent / 100)
```
