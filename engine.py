user_balances = {
    'BTC': 5.0,
    'ETH': 10.0,
    'USDT': 1000.0
}

order_book = {
    'BTC-ETH': {
        'bids': [{'price': 0.03, 'amount': 10}, {'price': 0.029, 'amount': 5}],
        'asks': [{'price': 0.032, 'amount': 7}, {'price': 0.033, 'amount': 12}]
    }
}

def execute_order(order_type, base_currency, quote_currency, amount, price=None, limit_price=None, stop_price=None):
    order_type_logic = {
        'buy': {
            'order_book_key': f'{base_currency}-{quote_currency}',
            'order_book_bid': -1,
            'order_book_ask_price_key': 'price',
            'balance_key': quote_currency,
            'other_balance_key': base_currency
        },
        'sell': {
            'order_book_key': f'{quote_currency}-{base_currency}',
            'order_book_bid': 0,
            'order_book_ask_price_key': 'price_inverse',
            'balance_key': base_currency,
            'other_balance_key': quote_currency
        }
    }

    logic = order_type_logic[order_type]

    if order_type == 'limit':
        if limit_price is None:
            print('Limit price must be specified for a limit order')
            return
        elif price is None:
            print('Price must be specified for a limit order')
            return
        else:
            if order_type == 'buy' and price >= limit_price:
                print('Limit price must be greater than the current price for a buy limit order')
                return
            elif order_type == 'sell' and price <= limit_price:
                print('Limit price must be less than the current price for a sell limit order')
                return
    elif order_type == 'stop':
        if stop_price is None:
            print('Stop price must be specified for a stop order')
            return
        elif price is None:
            print('Price must be specified for a stop order')
            return
        else:
            if order_type == 'buy' and price <= stop_price:
                print('Stop price must be less than the current price for a buy stop order')
                return
            elif order_type == 'sell' and price >= stop_price:
                print('Stop price must be greater than the current price for a sell stop order')
                return

    if order_type == 'market':
        if price is not None:
            print('Price should not be specified for a market order')
            return
        price = order_book[logic['order_book_key']]['asks' if order_type == 'buy' else 'bids'][0]['price']

    order_executed = False
    while amount > 0:
        if order_type in ['limit', 'market']:
            matching_orders = [order for order in order_book[logic['order_book_key']]['asks' if order_type == 'buy' else 'bids'] if order['price'] <= price]
        elif order_type == 'stop':
            matching_orders = [order for order in order_book[logic['order_book_key']]['asks' if order_type == 'buy' else 'bids'] if order['price'] >= stop_price]
        if len(matching_orders) == 0:
            break

        best_matching_order = matching_orders[0]
        order_amount = min(amount, best_matching_order['amount'])
        order_value = order_amount * best_matching_order[logic['order_book_ask_price_key']]

        # Update the user's balance
        user_balances[logic['balance_key']] -= order_value
        user_balances[logic['other_balance_key']] += order_amount

        # Update the order book
        best_matching_order['amount'] -= order_amount
        if best_matching_order['amount'] == 0:
            order_book[logic['order_book_key']]['asks' if order_type == 'buy' else 'bids'].remove(best_matching_order)

        # Set order_executed flag to True and break the loop
        order_executed = True
        break

    if order_executed:
        print(f'{order_type.capitalize()} order executed successfully')
    else:
        print(f'No matching orders found for {order_type} order')

    return user_balances, order_book
