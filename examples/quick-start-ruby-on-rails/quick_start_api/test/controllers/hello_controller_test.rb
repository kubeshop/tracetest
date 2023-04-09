require "test_helper"

class HelloControllerTest < ActionDispatch::IntegrationTest
  test "should get show" do
    get hello_show_url
    assert_response :success
  end
end
