Feature: 修改預設限額相關測試項目

  Background:
    Given 初始化用戶1限額、初始化用戶2限額

  Scenario: 用戶1『未客製化』限額時，更新預設值後用戶限額跟著『更新』
    Given 用戶1 DB 現在限額為：
      | daily_limit     | 300000  |
      | monthly_limit   | 1000000 |
      | is_customized   | false     |
    When 呼叫 GET /api/risk_control/8e8f09c4-aa2b-4b22-908b-28ec106641d6/role 取得該用戶設定
    Then 回傳的結果應為預設限額且包含：
      | daily_limit     | 300000  |
      | monthly_limit   | 1000000 |
    When 管理者呼叫 PATCH /api/config/limit API 並傳入以下參數：
      | role            | 1         |
      | daily_limit     | 1111      |
      | monthly_limit   | 2222      |
      | level1_volumn   | 500000    |
      | level2_volumn   | 2000000   |
      | level1_days     | 7         |
      | level2_days     | 60        |
      | velocity_days   | 1         |
      | velocity_times  | 5         |
      | reason          | test      |
    Then API 回傳的結果應為：
      """
      {
          "code": 20000,
          "msg": "success",
          "reason": "",
          "data": null
      }
      """
    When 呼叫 GET /api/risk_control/8e8f09c4-aa2b-4b22-908b-28ec106641d6/role 取得該用戶設定
    Then 用戶1 回傳的結果應包含：
      | daily_limit     | 1111    |
      | monthly_limit   | 2222    |
      | level1          | 500000  |
      | level2          | 2000000 |
      | level1_days     | 7       |
      | level2_days     | 60      |
      | velocity_days   | 1       |
      | velocity_times  | 5       |
    And 移除暫時DB

  # Scenario: 用戶2『客製化』限額時，更新預設值後用戶限額『不變』
  #   Given 呼叫 GET /api/risk_control/8e8f09c4-aa2b-4b22-908b-28ec106641d6/role 取得該用戶設定
  #   Then 回傳的結果應為預設限額且包含：
  #     | daily_limit     | 300000  |
  #     | monthly_limit   | 1000000 |
  #   When 呼叫 PATCH /api/risk_control/8e8f09c4-aa2b-4b22-908b-28ec106641d6/limit 修改用戶限額，並傳入以下參數：
  #     | daily_limit     | 1 |
  #     | monthly_limit   | 2 |
  #     | reason          | test      |
  #     | level1          | 3  |
  #     | level2          | 4 |
  #     | level1_days     | 5       |
  #     | level2_days     | 6      |
  #   Then API 回傳的結果應為：
  #     """
  #     {
  #         "code": 20000,
  #         "msg": "success",
  #         "reason": "",
  #         "data": null
  #     }
  #     """
  #   And 用戶2 DB 現在限額為：
  #     | daily_limit     | 1         |
  #     | monthly_limit   | 2         |
  #     | is_customized   | true      |
  #   When 管理者呼叫 PATCH /api/config/limit API 並傳入以下參數：
  #     | role            | 1         |
  #     | daily_limit     | 1111      |
  #     | monthly_limit   | 2222      |
  #     | level1_volumn   | 500000    |
  #     | level2_volumn   | 2000000   |
  #     | level1_days     | 7         |
  #     | level2_days     | 60        |
  #     | velocity_days   | 1         |
  #     | velocity_times  | 5         |
  #     | reason          | test      |
  #   Then API 回傳的結果應為：
  #     """
  #     {
  #         "code": 20000,
  #         "msg": "success",
  #         "reason": "",
  #         "data": null
  #     }
  #     """
  #   When 呼叫 GET /api/risk_control/8e8f09c4-aa2b-4b22-908b-28ec106641d6/role 取得該用戶設定
  #   Then 用戶2 回傳的結果應包含：
  #     | daily_limit     | 1         |
  #     | monthly_limit   | 2         |
  #     | level1          | 3         |
  #     | level2          | 4         |
  #     | level1_days     | 5         |
  #     | level2_days     | 6         |
  #     | velocity_days   | 1         |
  #     | velocity_times  | 5         |
  #   And 移除暫時DB

