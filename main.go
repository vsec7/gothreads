package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {

	switch mode := Mode; mode {
	case "login":
		login(Username, Password)
	case "like":
		for true {
			likeFeed()
			fmt.Printf("----------[ Delay for %d Seconds]----------\n", Delay)
			time.Sleep(time.Duration(Delay) * time.Second)
		}

	}
}

var (
	Username string
	Password string
	Delay    int
	Mode     string
)

func init() {
	flag.StringVar(&Username, "u", "", "Username")
	flag.StringVar(&Password, "p", "", "Password")
	flag.IntVar(&Delay, "d", 120, "Delay *Seconds")
	flag.StringVar(&Mode, "m", "", "Mode")

	flag.Usage = func() {
		h := []string{
			"",
			"GoThreads Bot",
			"",
			"Is a simple GO tool for threads.net",
			"",
			"Crafted By : github.com/vsec7",
			"",
			"Basic Usage :",
			" ▶ gothreads -m login -u username -p password",
			" ▶ gothreads -m like",
			"Options :",
			"  -m, --m <mode>	Set Mode: login, like",
			"  -d, --d <delay>	Set Delay *Seconds (default: 120)",
			"  -u, --u <username>	Set Instagram Username",
			"  -p, --p <password>	Set Instagram Password",
			"",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
	flag.Parse()
}

func likeFeed() {
	tokenBytes, err := ioutil.ReadFile("token.txt")
	if err != nil {
		fmt.Println("[x] Failed to read token.txt")
		fmt.Println("[▶] gothreads -m login -u username -p password")
		os.Exit(0)
	}

	token := strings.TrimSpace(string(tokenBytes))

	uuid := uuid()

	headers := map[string]string{
		"X-Bloks-Version-Id": "5f56efad68e1edec7801f630b5c122704ec5378adbee6609a448f105f34a9c73",
		"X-Ig-Www-Claim":     "hmac.AR2jr3_r-N6PqPM09G7tetqnPfD9P_Ux_HFjJvwyPwksRLqR",
		"X-Ig-Device-Id":     uuid,
		"X-Ig-Android-Id":    "android-6be35fa278d92525",
		"User-Agent":         "Barcelona 289.0.0.77.109 Android (31/12; 440dpi; 1080x2148; Google/google; sdk_gphone64_arm64; emulator64_arm64; ranchu; en_US; 489720145)",
		"Accept-Language":    "en-US",
		"Authorization":      "Bearer " + token,
		"Host":               "i.instagram.com",
	}

	getFeedsURL := "https://i.instagram.com/api/v1/feed/text_post_app_timeline/"
	postData := "feed_view_info=[]&max_id=&pagination_source=text_post_feed_threads&is_pull_to_refresh=0&_uuid=" + uuid + "&bloks_versioning_id=5f56efad68e1edec7801f630b5c122704ec5378adbee6609a448f105f34a9c73"

	getFeedsResp, err := request(getFeedsURL, postData, headers, true)
	if err != nil {
		fmt.Println("[x] Failed to get feeds:", err)
		return
	}

	parsegetFeeds := parseGetFeedsResponse(getFeedsResp)
	if parsegetFeeds == nil {
		fmt.Println("[x] Failed to parse feeds response.")
		return
	}

	if parsegetFeeds["message"] == "login_required" {
		fmt.Println("[x] Auth expired, please get token again.")
		fmt.Println("[▶] gothreads -m login -u username -p password")
		os.Exit(0)
	}

	items, ok := parsegetFeeds["items"].([]interface{})
	if !ok {
		fmt.Println("[x] Failed to extract items from feeds response.")
		return
	}

	for _, item := range items {
		postingan, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("[x] Failed to extract postingan from feeds response.")
			continue
		}

		threadItems, ok := postingan["thread_items"].([]interface{})
		if !ok {
			fmt.Println("[x] Failed to extract thread_items from postingan.")
			continue
		}

		if len(threadItems) > 0 {
			threadItem, ok := threadItems[0].(map[string]interface{})
			if !ok {
				fmt.Println("[x] Failed to extract threadItem from thread_items.")
				continue
			}

			post, ok := threadItem["post"].(map[string]interface{})
			if !ok {
				fmt.Println("[x] Failed to extract post from threadItem.")
				continue
			}

			mediaid, ok := post["id"].(string)
			if !ok {
				fmt.Println("[x] Failed to extract mediaid from post.")
				continue
			}

			fmt.Println("media_id:", mediaid)

			likeURL := "https://i.instagram.com/api/v1/media/" + mediaid + "/like/"
			likeData := "signed_body=SIGNATURE.%7B%22delivery_class%22%3A%22organic%22%2C%22tap_source%22%3A%22button%22%2C%22media_id%22%3A%22" + mediaid + "%22%2C%22radio_type%22%3A%22wifi-none%22%2C%22_uuid%22%3A%22" + uuid + "%2C%22recs_ix%22%3A%221%22%2C%22is_carousel_bumped_post%22%3A%22false%22%2C%22container_module%22%3A%22ig_text_feed_timeline%22%2C%22feed_position%22%3A%221%22%7D&d=0"

			likeResp, err := request(likeURL, likeData, headers, false)
			if err != nil {
				fmt.Println("[x] Failed to like post:", err)
				continue
			}

			fmt.Println(likeResp) // Status "ok" or "fail"
		}
	}
}

func login(username string, password string) {
	uuid := uuid()

	headers := map[string]string{
		"X-Bloks-Version-Id": "5f56efad68e1edec7801f630b5c122704ec5378adbee6609a448f105f34a9c73",
		"X-Ig-Www-Claim":     "hmac.AR2jr3_r-N6PqPM09G7tetqnPfD9P_Ux_HFjJvwyPwksRLqR",
		"X-Ig-Device-Id":     uuid,
		"X-Ig-Android-Id":    "android-6be35fa278d92525",
		"User-Agent":         "Barcelona 289.0.0.77.109 Android (31/12; 440dpi; 1080x2148; Google/google; sdk_gphone64_arm64; emulator64_arm64; ranchu; en_US; 489720145)",
		"Accept-Language":    "en-US",
		"Host":               "i.instagram.com",
	}

	getTokenURL := "https://i.instagram.com/api/v1/bloks/apps/com.bloks.www.bloks.caa.login.async.send_login_request/"
	postData := fmt.Sprintf(`params={"client_input_params":{"device_id":"android-6be35fa278d92525","login_attempt_count":1,"secure_family_device_id":"","machine_id":"ZKoroAABAAGyzN_tN5j5gN3Q0kpR","accounts_list":[],"auth_secure_device_id":"","password":"%s","family_device_id":"621af360-a821-4229-9e3a-678d59eb7d37","fb_ig_device_id":[],"device_emails":[],"try_num":1,"event_flow":"login_manual","event_step":"home_page","openid_tokens":{},"client_known_key_hash":"","contact_point":"%s","encrypted_msisdn":""},"server_params":{"username_text_input_id":"wktbih:48","device_id":"android-6be35fa278d92525","should_trigger_override_login_success_action":0,"server_login_source":"login","waterfall_id":"6b2356be-c2a1-41ac-8f81-9af5ec0aee87","login_source":"Login","INTERNAL__latency_qpl_instance_id":196987789700164,"is_platform_login":0,"credential_type":"password","family_device_id":"621af360-a821-4229-9e3a-678d59eb7d37","INTERNAL__latency_qpl_marker_id":36707139,"offline_experiment_group":"caa_iteration_v3_perf_ig_4","INTERNAL_INFRA_THEME":"harm_f","password_text_input_id":"wktbih:49","ar_event_source":"login_home_page"}}`,
		password, username)

	resp, err := request(getTokenURL, postData, headers, true)
	if err != nil {
		fmt.Println("[x] Failed to get auth token. Try again.\n[▶] gothreads -m login -u username -p password")
		os.Exit(0)
	}

	authToken := extractAuthToken(resp)
	if authToken == "" {
		fmt.Println("[x] Failed to get auth token. Try again.\n[▶] gothreads -m login -u username -p password")
		os.Exit(0)
	}

	err = ioutil.WriteFile("token.txt", []byte(authToken), 0644)
	if err != nil {
		fmt.Println("[x] Failed to write token to file:", err)
		os.Exit(0)
	}

	fmt.Println("[+] Successfully Login\n[+] Bearer Token Saved to token.txt")
}

func uuid() string {
	return fmt.Sprintf("%04x%04x-%04x-%04x-%04x-%04x%04x%04x",
		rand.Intn(0xffff),
		rand.Intn(0xffff),
		rand.Intn(0xffff),
		rand.Intn(0x0fff)|0x4000,
		rand.Intn(0x3fff)|0x8000,
		rand.Intn(0xffff),
		rand.Intn(0xffff),
		rand.Intn(0xffff),
	)
}

func parseGetFeedsResponse(response string) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		fmt.Println("Failed to parse Get Feeds response:", err)
		return nil
	}
	return result
}

func request(url, data string, headers map[string]string, outputHeader bool) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	if outputHeader {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractAuthToken(response string) string {
	re := regexp.MustCompile(`Bearer (.*?)\\\\`)
	matches := re.FindStringSubmatch(response)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}
