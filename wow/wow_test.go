package wow_test

import (
    "log"
    "os"
    "testing"
    "github.com/nezorflame/discord-wow-bot/wow"
)

func createMembersList() wow.MembersList {
    var membersList wow.MembersList
    membersList = append(membersList, wow.GuildMember{
        Member: wow.Character{
            Name: "AAAAA",
            Level: 110,
            Spec: wow.Specialization{
                Name: "Spec1",
            },
            Class: "Class1",
            Items: wow.Items{
                AvgItemLvlEq: 845,
            },
        },
        Rank: 1,
    })
    membersList = append(membersList, wow.GuildMember{
        Member: wow.Character{
            Name: "BBBB",
            Level: 110,
            Spec: wow.Specialization{
                Name: "Spec2",
            },
            Class: "Class1",
            Items: wow.Items{
                AvgItemLvlEq: 845,
            },
        },
        Rank: 1,
    })
    membersList = append(membersList, wow.GuildMember{
        Member: wow.Character{
            Name: "CCC",
            Level: 110,
            Spec: wow.Specialization{
                Name: "Spec2",
            },
            Class: "Class1",
            Items: wow.Items{
                AvgItemLvlEq: 847,
            },
        },
        Rank: 1,
    })
    membersList = append(membersList, wow.GuildMember{
        Member: wow.Character{
            Name: "DD",
            Level: 109,
            Spec: wow.Specialization{
                Name: "Spec1",
            },
            Class: "Class2",
            Items: wow.Items{
                AvgItemLvlEq: 847,
            },
        },
        Rank: 1,
    })
    membersList = append(membersList, wow.GuildMember{
        Member: wow.Character{
            Name: "E",
            Level: 109,
            Spec: wow.Specialization{
                Name: "Spec1",
            },
            Class: "Class3",
            Items: wow.Items{
                AvgItemLvlEq: 740,
            },
        },
        Rank: 1,
    })
    return membersList
}

func TestStringSorting(t *testing.T) {
    wow.Logger = log.New(os.Stderr, "  ", log.Ldate|log.Ltime)
    members := createMembersList()
    newMembers := members.SortGuildMembers([]string{"name=desc", "level=desc", "ilvl=desc"})
    // for _, m := range newMembers {
    //     log.Println(m.Member.Name, m.Member.Level, m.Member.Class, m.Member.Spec, m.Member.Items.AvgItemLvlEq)
    // }
    if len(newMembers) != len(members) {
        log.Println(newMembers)
        t.Fail()
    }
}