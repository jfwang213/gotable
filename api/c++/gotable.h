// Copyright 2015 stevejiang. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef _GO_TABLE_H_
#define _GO_TABLE_H_

#include <stdint.h>
#include <stdlib.h>
#include <string>
#include <set>
#include <vector>

namespace gotable {

using std::string;
using std::vector;

// GoTable Error Code List
enum {
	EcNotExist    = 1,   // Key NOT exist
	EcOk          = 0,   // Success
	EcCasNotMatch = -50, // CAS not match, get new CAS and try again
	EcTempFail    = -51, // Temporary failed, retry may fix this
	EcUnknownCmd  = -60, // Unknown cmd
	EcAuthFailed  = -61, // Authorize failed
	EcNoPrivilege = -62, // No access privilege
	EcWriteSlaver = -63, // Can NOT write slaver directly
	EcSlaverCas   = -64, // Invalid CAS on slaver for GET/MGET (cannot be 0)
	EcReadFail    = -65, // Read failed
	EcWriteFail   = -66, // Write failed
	EcDecodeFail  = -67, // Decode request PKG failed
	EcInvDbId     = -68, // Invalid DB ID (cannot be 255)
	EcInvRowKey   = -69, // RowKey length should be [1 ~ 255]
	EcInvValue    = -70, // Value length should be [0 ~ 1MB]
	EcInvPkgLen   = -71, // Pkg length should be less than 2MB
	EcInvScanNum  = -72, // Scan request number out of range
	EcScanEnded   = -73, // Already scan/dump to end
};

struct GetArgs {
	uint8_t tableId;
	string  rowKey;
	string  colKey;
	uint32_t cas;

	GetArgs() : tableId(0), cas(0) {}

	GetArgs(uint8_t tableId, const string& rowKey, const string& colKey, uint32_t cas) :
			tableId(tableId), rowKey(rowKey), colKey(colKey), cas(cas) {}
};

struct GetReply {
	int8_t  errCode; // Error Code Replied
	uint8_t tableId;
	string  rowKey;
	string  colKey;
	string  value;
	int64_t score;
	uint32_t cas;

	GetReply() : errCode(0), tableId(0), score(0), cas(0) {}
};

struct SetArgs {
	uint8_t tableId;
	string  rowKey;
	string  colKey;
	string  value;
	int64_t score;
	uint32_t cas;

	SetArgs() : tableId(0), score(0), cas(0) {}

	SetArgs(uint8_t tableId, const string& rowKey, const string& colKey,
			const string& value, int64_t score, uint32_t cas) :
			tableId(tableId), rowKey(rowKey), colKey(colKey),
			value(value), score(score), cas(cas) {}
};

struct SetReply {
	int8_t  errCode; // Error Code Replied
	uint8_t tableId;
	string  rowKey;
	string  colKey;

	SetReply() : errCode(0), tableId(0) {}
};

struct IncrArgs {
	uint8_t tableId;
	string  rowKey;
	string  colKey;
	int64_t score;
	uint32_t cas;

	IncrArgs() : tableId(0), score(0), cas(0) {}

	IncrArgs(uint8_t tableId, const string& rowKey, const string& colKey,
			int64_t score, uint32_t cas) :
			tableId(tableId), rowKey(rowKey), colKey(colKey), score(score), cas(cas) {}
};

struct IncrReply {
	int8_t  errCode; // Error Code Replied
	uint8_t tableId;
	string  rowKey;
	string  colKey;
	string  value;
	int64_t score;

	IncrReply() : errCode(0), tableId(0), score(0) {}
};

typedef GetArgs DelArgs;
typedef SetReply DelReply;

struct ScanKV {
	string  colKey;
	string  value;
	int64_t score;

	ScanKV() : score(0) {}
};

struct ScanReply {
	uint8_t tableId;
	string  rowKey;
	vector<ScanKV> kvs;
	bool    end;    // true: Scan to end, stop now

	ScanReply() : tableId(0), kvs(), end(false) {}

private:
	struct ScanContext {
		bool zop;
		bool asc;          // true: Ascending  order; false: Descending  order
		bool orderByScore; // true: Score+ColKey; false: ColKey
		int num;           // Max number of scan reply records
	};
	ScanContext ctx;
	friend class Client;
};

struct DumpKV {
	uint8_t tableId;
	uint8_t colSpace;
	string  rowKey;
	string  colKey;
	string  value;
	int64_t score;

	DumpKV() : tableId(0), colSpace(0), score(0) {}
};

struct DumpReply {
	vector<DumpKV> kvs;
	bool end;    // true: Dump to end, stop now

private:
	struct DumpContext {
		bool     oneTable;     // Never change during dump
		uint8_t  tableId;      // Never change during dump
		uint16_t startUnitId;  // Never change during dump
		uint16_t endUnitId;    // Never change during dump
		uint16_t lastUnitId;   // The last unit ID tried to dump
		bool     unitStart;    // Next dump start from new UnitId
	};
	DumpContext ctx;
	friend class Client;
};

struct PkgOneOp;
struct PkgMultiOp;
struct PkgDumpResp;

class Client {
public:
	Client(int fd);
	~Client();

	static Client* Dial(const char* ip, int port);
	void close();

	void select(uint8_t dbId);
	uint8_t databaseId();
	int auth(const char* password);
	int ping();

	int get(uint8_t tableId, const string& rowKey, const string& colKey,
			string* value, int64_t* score, uint32_t* cas=NULL);
	int zGet(uint8_t tableId, const string& rowKey, const string& colKey,
			string* value, int64_t* score, uint32_t* cas=NULL);
	int set(uint8_t tableId, const string& rowKey, const string& colKey,
			const string& value, int64_t score, uint32_t cas=0);
	int zSet(uint8_t tableId, const string& rowKey, const string& colKey,
			const string& value, int64_t score, uint32_t cas=0);
	int del(uint8_t tableId, const string& rowKey, const string& colKey,
			uint32_t cas=0);
	int zDel(uint8_t tableId, const string& rowKey, const string& colKey,
			uint32_t cas=0);
	int incr(uint8_t tableId, const string& rowKey, const string& colKey,
			string* value, int64_t* score, uint32_t cas=0);
	int zIncr(uint8_t tableId, const string& rowKey, const string& colKey,
			string* value, int64_t* score, uint32_t cas=0);

	int mGet(const vector<GetArgs>& args, vector<GetReply>* reply);
	int zmGet(const vector<GetArgs>& args, vector<GetReply>* reply);
	int mSet(const vector<SetArgs>& args, vector<SetReply>* reply);
	int zmSet(const vector<SetArgs>& args, vector<SetReply>* reply);
	int mDel(const vector<DelArgs>& args, vector<DelReply>* reply);
	int zmDel(const vector<DelArgs>& args, vector<DelReply>* reply);
	int mIncr(const vector<IncrArgs>& args, vector<IncrReply>* reply);
	int zmIncr(const vector<IncrArgs>& args, vector<IncrReply>* reply);

	int scan(uint8_t tableId, const string& rowKey, const string& colKey,
			bool asc, int num, ScanReply* reply);
	int scanStart(uint8_t tableId, const string& rowKey,
			bool asc, int num, ScanReply* reply);
	int zScan(uint8_t tableId, const string& rowKey, const string& colKey, int64_t score,
			bool asc, bool orderByScore, int num, ScanReply* reply);
	int zScanStart(uint8_t tableId, const string& rowKey,
			bool asc, bool orderByScore, int num, ScanReply* reply);
	int scanMore(const ScanReply& last, ScanReply* reply);

	int dump(bool oneTable, uint8_t tableId, uint8_t colSpace,
			const string& rowKey, const string& colKey, int64_t score,
			uint16_t startUnitId, uint16_t endUnitId, DumpReply* reply);
	int dumpDB(DumpReply* reply);
	int dumpTable(uint8_t tableId, DumpReply* reply);
	int dumpMore(const DumpReply& last, DumpReply* reply);

private:
	int doOneOp(bool zop, uint8_t cmd, uint8_t tableId,
			const string& rowKey, const string& colKey,
			const string& value, int64_t score, uint32_t cas,
			PkgOneOp* reply, string& pkg);

	template <typename T>
	int doMultiOp(bool zop, uint8_t cmd, const vector<T>& args,
			PkgMultiOp* reply, string& pkg);

	int doScan(bool zop, uint8_t tableId, const string& rowKey, const string& colKey,
			int64_t score, bool start, bool asc, bool orderByScore, int num,
			ScanReply* reply, PkgMultiOp* resp, string& pkg);

	int doDump(bool oneTable, uint8_t tableId, uint8_t colSpace,
			const string& rowKey, const string& colKey, int64_t score,
			uint16_t startUnitId, uint16_t endUnitId,
			DumpReply* reply, PkgDumpResp* resp, string& pkg);

private:  //disable
	Client(const Client&);
	void operator=(const Client&);

private:
	bool     closed;
	int      fd;
	uint8_t  dbId;
	uint64_t seq;
	bool              authAdmin;
	std::set<uint8_t> setAuth;
	char     buf[4096];
};

}  // namespace gotable
#endif
